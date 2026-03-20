import os
import json
import zipfile
import base64
import hashlib
from pathlib import Path
from .crypto_utils import CryptoUtils


class ZIPBuilder:
    def __init__(self, template_dir: str = None):
        self.template_dir = template_dir or os.path.join(os.path.dirname(__file__), "templates")

    def build(
        self,
        package_id: str,
        package_name: str,
        source_file: str,
        output_path: str,
        verify_config: dict,
        private_key: bytes = None,
        aes_key: bytes = None
    ) -> str:
        if not os.path.exists(source_file):
            raise FileNotFoundError(f"Source file not found: {source_file}")
        
        os.makedirs(os.path.dirname(output_path), exist_ok=True)
        
        if private_key is None or aes_key is None:
            private_key, public_key = CryptoUtils.generate_rsa_keypair(2048)
            aes_key = CryptoUtils.generate_aes_key(256)
        else:
            from Crypto.PublicKey import RSA
            public_key = RSA.import_key(private_key).publickey().export_key()
        
        with open(source_file, 'rb') as f:
            source_data = f.read()
        encrypted_data = CryptoUtils.encrypt_aes_with_key(source_data, aes_key)
        
        password_list = verify_config.get("passwords", [])
        encrypted_keys_per_password = {}
        password_hints = []
        
        for pwd_info in password_list:
            password_key_hash = pwd_info.get("password")
            if password_key_hash:
                derived_key, salt = CryptoUtils.derive_key_from_password(password_key_hash)
                encrypted_key = CryptoUtils.encrypt_aes_with_key(aes_key, derived_key)
                salt_b64 = base64.b64encode(salt).decode()
                encrypted_keys_per_password[password_key_hash] = {
                    "salt": salt_b64,
                    "encrypted_key": encrypted_key
                }
                password_hints.append({
                    "priority": pwd_info.get("priority"),
                    "valid_from": pwd_info.get("valid_from"),
                    "valid_until": pwd_info.get("valid_until"),
                    "password_hint": password_key_hash[:8]
                })
        
        config_data = {
            "version": "1.0.0",
            "package_id": package_id,
            "package_name": package_name,
            "encrypted_keys_per_password": encrypted_keys_per_password,
            "password_hints": password_hints,
        }
        
        source_hash = CryptoUtils.compute_file_hash(source_file)
        
        try:
            import pyzipper
            with pyzipper.AESZipFile(output_path, 'w', compression=pyzipper.ZIP_DEFLATED, encryption=pyzipper.WZ_AES) as zf:
                zf.setpassword(b'PTM_PACKAGE_KEY_2024')
                
                temp_encrypted = os.path.join(os.path.dirname(output_path) or '.', ".encrypted.tmp")
                with open(temp_encrypted, 'wb') as ef:
                    ef.write(encrypted_data.encode('utf-8'))
                zf.write(temp_encrypted, arcname="data.bin.encrypted")
                os.remove(temp_encrypted)
                
                config_path = os.path.join(os.path.dirname(output_path) or '.', ".config.tmp")
                with open(config_path, 'w', encoding='utf-8') as f:
                    json.dump(config_data, f, indent=2, ensure_ascii=False)
                zf.write(config_path, arcname="config.json")
                os.remove(config_path)
                
                info_data = {
                    "package_id": package_id,
                    "package_name": package_name,
                    "source_file_hash": source_hash,
                    "format": "zip"
                }
                info_path = os.path.join(os.path.dirname(output_path) or '.', ".info.tmp")
                with open(info_path, 'w', encoding='utf-8') as f:
                    json.dump(info_data, f, indent=2)
                zf.write(info_path, arcname="info.json")
                os.remove(info_path)
        except ImportError:
            with zipfile.ZipFile(output_path, 'w', zipfile.ZIP_DEFLATED) as zf:
                temp_encrypted = os.path.join(os.path.dirname(output_path) or '.', ".encrypted.tmp")
                with open(temp_encrypted, 'w', encoding='utf-8') as ef:
                    ef.write(encrypted_data)
                zf.write(temp_encrypted, arcname="data.bin.encrypted")
                os.remove(temp_encrypted)
                
                config_path = os.path.join(os.path.dirname(output_path) or '.', ".config.tmp")
                with open(config_path, 'w', encoding='utf-8') as f:
                    json.dump(config_data, f, indent=2, ensure_ascii=False)
                zf.write(config_path, arcname="config.json")
                os.remove(config_path)
                
                info_data = {
                    "package_id": package_id,
                    "package_name": package_name,
                    "source_file_hash": source_hash,
                    "format": "zip"
                }
                info_path = os.path.join(os.path.dirname(output_path) or '.', ".info.tmp")
                with open(info_path, 'w', encoding='utf-8') as f:
                    json.dump(info_data, f, indent=2)
                zf.write(info_path, arcname="info.json")
                os.remove(info_path)
        
        return output_path

    def extract(self, zip_path: str, extract_dir: str) -> dict:
        if not os.path.exists(zip_path):
            raise FileNotFoundError(f"ZIP file not found: {zip_path}")
        
        os.makedirs(extract_dir, exist_ok=True)
        
        with zipfile.ZipFile(zip_path, 'r') as zf:
            zf.extractall(extract_dir)
            
            info_path = os.path.join(extract_dir, "info.json")
            if os.path.exists(info_path):
                with open(info_path, 'r', encoding='utf-8') as f:
                    info = json.load(f)
                return info
            
            config_path = os.path.join(extract_dir, "config.json")
            if os.path.exists(config_path):
                with open(config_path, 'r', encoding='utf-8') as f:
                    config = json.load(f)
                return config
        
        return {}
    
    def extract_with_password(self, zip_path: str, extract_dir: str, password: bytes) -> dict:
        if not os.path.exists(zip_path):
            raise FileNotFoundError(f"ZIP file not found: {zip_path}")
        
        os.makedirs(extract_dir, exist_ok=True)
        
        import pyzipper
        with pyzipper.AESZipFile(zip_path, 'r') as zf:
            zf.setpassword(password)
            zf.extractall(extract_dir)
            
            info_path = os.path.join(extract_dir, "info.json")
            if os.path.exists(info_path):
                with open(info_path, 'r', encoding='utf-8') as f:
                    info = json.load(f)
                return info
        
        return {}
    
    @staticmethod
    def verify_and_decrypt(config: dict, password: str) -> bytes:
        password_key_hash = hashlib.sha256(password.encode('utf-8')).hexdigest()[:16]
        encrypted_keys = config.get("encrypted_keys_per_password", {})
        
        key_data = encrypted_keys.get(password_key_hash)
        if not key_data:
            raise ValueError("Invalid password")
        
        salt = base64.b64decode(key_data["salt"])
        derived_key, _ = CryptoUtils.derive_key_from_password(password_key_hash, salt)
        aes_key = CryptoUtils.decrypt_aes_with_key(key_data["encrypted_key"], derived_key)
        return aes_key
