# Copyright 2014 ETH Zurich
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""
:mod:`trc` --- SCION TRC parser
===============================================
"""
# Stdlib
import base64
import copy
import json
import os
import time

# External
import lz4

# SCION
from lib.crypto.asymcrypto import verify, sign
from lib.crypto.util import CERT_DIR
from lib.errors import SCIONParseError, SCIONVerificationError
from lib.packet.scion_addr import ISD_AS

ISD_STRING = 'ISD'
DESCRIPTION_STRING = 'Description'
VERSION_STRING = 'Version'
CREATION_TIME_STRING = 'CreationTime'
EXPIRATION_TIME_STRING = 'ExpirationTime'
CORE_ASES_STRING = 'CoreASes'
ROOT_CAS_STRING = 'RootCAs'
CERT_LOGS_STRING = 'CertLogs'
THRESHOLD_EEPKI_STRING = 'ThresholdEEPKI'
RAINS_STRING = 'RAINS'
QUORUM_TRC_STRING = 'QuorumTRC'
QUORUM_CAS_STRING = 'QuorumCAs'
GRACE_PERIOD_STRING = 'GracePeriod'
QUARANTINE_STRING = 'Quarantine'
SIGNATURES_STRING = 'Signatures'

ARPKI_KEY_STRING = 'ARPKIKey'
ARPKI_SRV_STRING = 'ARPKISrv'
CERTIFICATE_STRING = 'Certificate'
OFFLINE_KEY_ALG_STRING = 'OfflineKeyAlg'
OFFLINE_KEY_STRING = 'OfflineKey'
ONLINE_KEY_ALG_STRING = 'OnlineKeyAlg'
ONLINE_KEY_STRING = 'OnlineKey'
ROOT_RAINS_KEY_STRING = 'RootRAINSKey'
TRC_SRV_STRING = 'TRCSrv'


def get_trc_file_path(conf_dir, isd, version):  # pragma: no cover
    """
    Return the TRC file path for a given ISD.
    """
    return os.path.join(conf_dir, CERT_DIR, 'ISD%s-V%s.trc' % (isd, version))


class TRC(object):
    """
    The TRC class parses the TRC file of an ISD and stores such
    information for further use.

    :ivar int isd: the ISD identifier.
    :ivar str description: is a human readable description of an ISD.
    :ivar int version: the TRC file version.
    :ivar int create_time: the TRC file creation timestamp.
    :ivar int exp_time: the TRC expiration timestamp.
    :ivar dict core_ases: the set of core ASes and their certificates.
    :ivar dict root_cas: the set of root CAs and their certificates.
    :ivar dict cert_logs: is a dictionary of end entity certificate log servers of
        form {name: {"isd_as IP": pub_key}}
    :ivar int threshold_eepki: is a threshold number (nonnegative integer) of
        CAs that have to sign a domain’s policy
    :ivar dict rains: the RAINS section.
    :ivar int quorum_trc: number of core ASes necessary to sign a new TRC.
    :ivar int quorum_cas: number of CAs necessary to change CA entries
    :ivar int grace_period: defines for how long this TRC is valid when a new
        TRC is available
    :ivar bool quarantine: flag defining whether TRC is valid(quarantine=false)
        or an early announcement(quarantine=true)
    :ivar dict signatures: signatures generated by a quorum of trust roots.
    """

    VALIDITY_PERIOD = 365 * 24 * 60 * 60

    FIELDS_MAP = {
        ISD_STRING: ("isd", int),
        DESCRIPTION_STRING: ("description", str),
        VERSION_STRING: ("version", int),
        CREATION_TIME_STRING: ("create_time", int),
        EXPIRATION_TIME_STRING: ("exp_time", int),
        CORE_ASES_STRING: ("core_ases", dict),
        ROOT_CAS_STRING: ("root_cas", dict),
        CERT_LOGS_STRING: ("cert_logs", dict),
        THRESHOLD_EEPKI_STRING: ("threshold_eepki", int),
        RAINS_STRING: ("rains", dict),
        QUORUM_TRC_STRING: ("quorum_trc", int),
        QUORUM_CAS_STRING: ("quorum_cas", int),
        QUARANTINE_STRING: ("quarantine", bool),
        SIGNATURES_STRING: ("signatures", dict),
        GRACE_PERIOD_STRING: ("grace_period", int),
    }

    # list of fields in a dict of dicts which have to be encoded/decoded from base64
    MULTI_DICT_DECODE_FIELDS = {
        CORE_ASES_STRING: [ONLINE_KEY_STRING, OFFLINE_KEY_STRING],
        ROOT_CAS_STRING: [CERTIFICATE_STRING, ONLINE_KEY_STRING, ARPKI_KEY_STRING],
    }

    # list of fields in a dict which have to be encoded/decoded
    SIMPLE_DICT_DECODE_FIELDS = {
        RAINS_STRING: [ROOT_RAINS_KEY_STRING, ONLINE_KEY_STRING],
        SIGNATURES_STRING: [],
    }

    def __init__(self, trc_dict):
        """
        :param dict trc_dict: TRC as dict.
        """
        for k, (name, type_) in self.FIELDS_MAP.items():
            val = trc_dict[k]
            if type_ in (int,):
                val = int(val)
            elif type_ in (dict, ):
                val = copy.deepcopy(val)
            setattr(self, name, val)

        for attr, decode_list in self.MULTI_DICT_DECODE_FIELDS.items():
            field = getattr(self, self.FIELDS_MAP[attr][0])
            for entry in field.values():
                for key in decode_list:
                    entry[key] = base64.b64decode(entry[key].encode('utf-8'))

        for attr, decode_list in self.SIMPLE_DICT_DECODE_FIELDS.items():
            entry = getattr(self, self.FIELDS_MAP[attr][0])
            if not entry:
                continue
            for key in decode_list or entry:
                entry[key] = base64.b64decode(entry[key].encode('utf-8'))

        for subject, entry in trc_dict[CERT_LOGS_STRING].items():
            try:
                addr, pub_key = next(iter(entry.items()))
                self.cert_logs[subject][addr] = base64.b64decode(pub_key.encode('utf-8'))
            except StopIteration:
                raise SCIONParseError("Invalid CertLogs entry for %s: %s", subject, entry)

    def get_isd_ver(self):
        return self.isd, self.version

    def get_core_ases(self):
        res = []
        for key in self.core_ases:
            res.append(ISD_AS(key))
        return res

    def dict(self, with_signatures):
        """
        Return the TRC information.

        :param bool with_signatures:
            If True, include signatures in the return value.
        :returns: the TRC information.
        :rtype: dict
        """
        trc_dict = {}
        for k, (name, _) in self.FIELDS_MAP.items():
            trc_dict[k] = getattr(self, name)
        if not with_signatures:
            del trc_dict[SIGNATURES_STRING]
        return trc_dict

    @classmethod
    def from_raw(cls, trc_raw, lz4_=False):
        if lz4_:
            trc_raw = lz4.loads(trc_raw).decode("utf-8")
        trc = json.loads(trc_raw)
        return TRC(trc)

    @classmethod
    def from_values(cls, isd, description, version, core_ases, root_cas,
                    cert_logs, threshold_eepki, rains, quorum_trc,
                    quorum_cas, grace_period, quarantine, signatures, validity_period):
        """
        Generate a TRC instance.
        """
        now = int(time.time())
        trc_dict = {
            ISD_STRING: isd,
            DESCRIPTION_STRING: description,
            VERSION_STRING: version,
            CREATION_TIME_STRING: now,
            EXPIRATION_TIME_STRING: now + validity_period,
            CORE_ASES_STRING: core_ases,
            ROOT_CAS_STRING: root_cas,
            CERT_LOGS_STRING: cert_logs,
            THRESHOLD_EEPKI_STRING: threshold_eepki,
            RAINS_STRING: rains,
            QUORUM_TRC_STRING: quorum_trc,
            QUORUM_CAS_STRING: quorum_cas,
            GRACE_PERIOD_STRING: grace_period,
            QUARANTINE_STRING: quarantine,
            SIGNATURES_STRING: signatures,
        }
        trc = TRC(trc_dict)
        return trc

    def sign(self, isd_as, sig_priv_key):
        """
        Sign TRC and add computed signature to the TRC.

        :param ISD_AS isd_as: the ISD-AS of signer.
        :param SigningKey sig_priv_key: the signing key of signer.
        """
        data = self._sig_input()
        self.signatures[str(isd_as)] = sign(data, sig_priv_key)

    def _sig_input(self):
        d = self.dict(False)
        for k in d:
            if self.FIELDS_MAP[k][1] == dict:
                d[k] = self._encode_dict(d[k])
        j = json.dumps(d, sort_keys=True, separators=(',', ':'))
        return j.encode('utf-8')

    def _encode_dict(self, dict_):
        encoded_dict = {}
        for key, val in dict_.items():
            if type(val) is dict:
                val = self._encode_sub_dict(val)
            elif type(val) is bytes:
                val = base64.b64encode(val).decode('utf-8')
            encoded_dict[key] = val
        return encoded_dict

    def _encode_sub_dict(self, dict_):
        encoded_dict = {}
        for key, val in dict_.items():
            if type(val) is bytes:
                val = base64.b64encode(val).decode('utf-8')
            encoded_dict[key] = val
        return encoded_dict

    def to_json(self, with_signatures=True):
        """
        Convert the instance to json format.
        """
        trc_dict = copy.deepcopy(self.dict(with_signatures))
        for field, decode_list in self.MULTI_DICT_DECODE_FIELDS.items():
            for entry in trc_dict[field].values():
                for key in decode_list:
                    entry[key] = base64.b64encode(entry[key]).decode('utf-8')
        for field, decode_list in self.SIMPLE_DICT_DECODE_FIELDS.items():
            entry = trc_dict.get(field, None)
            if not entry or (field == SIGNATURES_STRING and not with_signatures):
                continue
            # Every value is decoded, if decode_list is empty
            for key in decode_list or entry:
                entry[key] = base64.b64encode(entry[key]).decode('utf-8')
        cert_logs = {}
        for subject, entry in trc_dict[CERT_LOGS_STRING].items():
            try:
                addr = next(iter(entry.keys()))
                entry[addr] = base64.b64encode(entry[addr]).decode('utf-8')
                cert_logs[subject] = entry
            except StopIteration:
                pass
        trc_dict[CERT_LOGS_STRING] = cert_logs
        trc_str = json.dumps(trc_dict, sort_keys=True, indent=4)
        return trc_str

    def pack(self, lz4_=False):
        ret = self.to_json().encode('utf-8')
        if lz4_:
            return lz4.dumps(ret)
        return ret

    def __str__(self):
        return self.to_json()

    def __eq__(self, other):  # pragma: no cover
        return str(self) == str(other)

    def check_active(self, max_trc=None):
        """
        Check if trusted TRC is active and can be used for certificate chain verification.

        :param TRC max_trc: newest available TRC for same ISD. (If none, self is newest TRC)
        :raises: SCIONVerificationError
        """
        now = int(time.time())
        if not (self.create_time <= now <= self.exp_time):
            raise SCIONVerificationError("Current time outside of validity period. "
                                         "Now %s Creation %s Expiration %s" %
                                         (now, self.create_time, self.exp_time))
        if not max_trc or self.version == max_trc.version:
            return
        if self.version + 1 != max_trc.version:
            raise SCIONVerificationError("Inactive TRC version: %s. Expected %s or %s" % (
                self.version, max_trc.version, max_trc.version - 1))
        if now > max_trc.create_time + max_trc.grace_period:
            raise SCIONVerificationError("TRC grace period has passed. Now %s Expiration %s" % (
                now, max_trc.create_time + max_trc.grace_period))

    def verify(self, trusted_trc):
        """
        Verify TRC based on a trusted TRC.

        :param TRC trusted_trc: a verified TRC, used as a trust anchor.
        :raises: SCIONVerificationError
        """
        if self.isd == trusted_trc.isd:
            self.verify_update(trusted_trc)
        else:
            self.verify_xsig(trusted_trc)

    def verify_update(self, old_trc):
        """
        Verify TRC update.
        Unsuccessful verification raises an error.

        :param TRC old_trc: a verified TRC, used as a trust anchor.
        :raises: SCIONVerificationError
        """
        if old_trc.isd != self.isd:
            raise SCIONVerificationError("Invalid TRC ISD %s. Expected %s" % (
                self.isd, old_trc.isd))
        if old_trc.version + 1 != self.version:
            raise SCIONVerificationError("Invalid TRC version %s. Expected %s" % (
                self.version, old_trc.isd + 1))
        if self.create_time < old_trc.create_time + old_trc.grace_period:
            raise SCIONVerificationError("Invalid timestamp %s. Expected > %s " % (
                self.create_time, old_trc.create_time + old_trc.grace_period))
        if self.quarantine or old_trc.quarantine:
            raise SCIONVerificationError("Early announcement")
        self._verify_signatures(old_trc)

    def verify_xsig(self, neigh_trc):
        """
        Verify cross signatures.

        :param TRC neigh_trc: neighbour TRC, used as a trust anchor.
        :raises: SCIONVerificationError
        """
        pass

    def _verify_signatures(self, old_trc):
        """
        Perform signature verification for core signatures as defined
        in old TRC. Raises an error if verification is unsuccessful.

        :param: TRC old_trc: the previous TRC which has already been verified.
        :raises: SCIONVerificationError
        """
        # Only look at signatures which are from core ASes as defined in old TRC
        val_count = 0
        # Count number of verifiable signatures
        for signer in old_trc.core_ases.keys():
            public_key = old_trc.core_ases[signer][ONLINE_KEY_STRING]
            try:
                verify(self._sig_input(), self.signatures[signer], public_key)
                val_count += 1
            except (SCIONVerificationError, KeyError):
                continue
        # Check if enough valid signatures
        if val_count < old_trc.quorum_trc:
            raise SCIONVerificationError("Not enough valid signatures %s. Expected %s" % (
                val_count, old_trc.quorum_trc))
