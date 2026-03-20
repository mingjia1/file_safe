import os
import json
import zipfile
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
        
        config_data = json.dumps(verify_config, ensure_ascii=False).encode('utf-8')
        config_hash = CryptoUtils.hash_sha256(config_data)
        
        encrypted_config = CryptoUtils.encrypt_aes_with_key(config_data, aes_key)
        encrypted_aes_key = CryptoUtils.encrypt_rsa(aes_key, public_key)
        
        package_info = CryptoUtils.create_config_package(
            package_id=package_id,
            package_name=package_name,
            public_key=public_key,
            encrypted_aes_key=encrypted_aes_key,
            encrypted_config=encrypted_config,
            config_hash=config_hash
        )
        
        source_hash = CryptoUtils.compute_file_hash(source_file)
        
        with zipfile.ZipFile(output_path, 'w', zipfile.ZIP_DEFLATED) as zf:
            zf.write(source_file, arcname="data.bin")
            
            config_path = os.path.join(os.path.dirname(output_path) or '.', "config.json")
            with open(config_path, 'w', encoding='utf-8') as f:
                json.dump(package_info, f, indent=2, ensure_ascii=False)
            zf.write(config_path, arcname="config.json")
            
            info_data = {
                "package_id": package_id,
                "package_name": package_name,
                "source_file_hash": source_hash,
                "format": "zip"
            }
            info_path = os.path.join(os.path.dirname(output_path) or '.', "info.json")
            with open(info_path, 'w', encoding='utf-8') as f:
                json.dump(info_data, f, indent=2)
            zf.write(info_path, arcname="info.json")
        
        if os.path.exists(config_path):
            os.remove(config_path)
        if os.path.exists(info_path):
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
