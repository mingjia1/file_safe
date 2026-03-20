import os
import hashlib
import json
import base64

try:
    from Crypto.Cipher import AES, PKCS1_OAEP
    from Crypto.PublicKey import RSA
    from Crypto.Random import get_random_bytes
    HAS_PYCRYPTODOME = True
except ImportError:
    HAS_PYCRYPTODOME = False


class CryptoUtils:
    @staticmethod
    def generate_aes_key(key_length: int = 256) -> bytes:
        if key_length == 256:
            return get_random_bytes(32)
        elif key_length == 192:
            return get_random_bytes(24)
        elif key_length == 128:
            return get_random_bytes(16)
        else:
            raise ValueError(f"Unsupported AES key length: {key_length}")

    @staticmethod
    def generate_rsa_keypair(key_length: int = 2048) -> tuple:
        if not HAS_PYCRYPTODOME:
            raise ImportError("pycryptodome is required for RSA operations")
        key = RSA.generate(key_length)
        private_key = key.export_key()
        public_key = key.publickey().export_key()
        return private_key, public_key

    @staticmethod
    def encrypt_aes(data: bytes, key: bytes) -> str:
        if not HAS_PYCRYPTODOME:
            raise ImportError("pycryptodome is required for AES operations")
        cipher = AES.new(key, AES.MODE_GCM)
        ciphertext, tag = cipher.encrypt_and_digest(data)
        return base64.b64encode(cipher.nonce + tag + ciphertext).decode()

    @staticmethod
    def decrypt_aes(encrypted_data: str, key: bytes) -> bytes:
        if not HAS_PYCRYPTODOME:
            raise ImportError("pycryptodome is required for AES operations")
        data = base64.b64decode(encrypted_data)
        nonce = data[:16]
        tag = data[16:32]
        ciphertext = data[32:]
        cipher = AES.new(key, AES.MODE_GCM, nonce=nonce)
        return cipher.decrypt_and_verify(ciphertext, tag)

    @staticmethod
    def encrypt_aes_with_key(data: bytes, key: bytes) -> str:
        if not HAS_PYCRYPTODOME:
            raise ImportError("pycryptodome is required for AES operations")
        cipher = AES.new(key, AES.MODE_EAX)
        ciphertext, tag = cipher.encrypt_and_digest(data)
        return base64.b64encode(cipher.nonce + tag + ciphertext).decode()

    @staticmethod
    def decrypt_aes_with_key(encrypted_data: str, key: bytes) -> bytes:
        if not HAS_PYCRYPTODOME:
            raise ImportError("pycryptodome is required for AES operations")
        data = base64.b64decode(encrypted_data)
        nonce = data[:16]
        tag = data[16:32]
        ciphertext = data[32:]
        cipher = AES.new(key, AES.MODE_EAX, nonce=nonce)
        return cipher.decrypt_and_verify(ciphertext, tag)

    @staticmethod
    def encrypt_rsa(data: bytes, public_key: bytes) -> str:
        if not HAS_PYCRYPTODOME:
            raise ImportError("pycryptodome is required for RSA operations")
        key = RSA.import_key(public_key)
        cipher = PKCS1_OAEP.new(key)
        return base64.b64encode(cipher.encrypt(data)).decode()

    @staticmethod
    def decrypt_rsa(encrypted_data: str, private_key: bytes) -> bytes:
        if not HAS_PYCRYPTODOME:
            raise ImportError("pycryptodome is required for RSA operations")
        key = RSA.import_key(private_key)
        cipher = PKCS1_OAEP.new(key)
        return cipher.decrypt(base64.b64decode(encrypted_data))

    @staticmethod
    def hash_sha256(data: bytes) -> str:
        return hashlib.sha256(data).hexdigest()

    @staticmethod
    def hash_sha512(data: bytes) -> str:
        return hashlib.sha512(data).hexdigest()

    @staticmethod
    def compute_file_hash(file_path: str, algorithm: str = "sha256") -> str:
        if algorithm == "sha256":
            hasher = hashlib.sha256()
        elif algorithm == "sha512":
            hasher = hashlib.sha512()
        else:
            raise ValueError(f"Unsupported hash algorithm: {algorithm}")
        
        with open(file_path, 'rb') as f:
            for chunk in iter(lambda: f.read(8192), b''):
                hasher.update(chunk)
        return hasher.hexdigest()

    @staticmethod
    def create_config_package(
        package_id: str,
        package_name: str,
        public_key: bytes,
        encrypted_aes_key: str,
        encrypted_config: str,
        config_hash: str
    ) -> dict:
        return {
            "version": "1.0.0",
            "package_id": package_id,
            "package_name": package_name,
            "public_key": base64.b64encode(public_key).decode() if isinstance(public_key, bytes) else public_key,
            "encrypted_aes_key": encrypted_aes_key,
            "encrypted_config": encrypted_config,
            "config_hash": config_hash,
        }

    @staticmethod
    def verify_config_integrity(config_data: bytes, expected_hash: str) -> bool:
        computed_hash = CryptoUtils.hash_sha256(config_data)
        return computed_hash == expected_hash
