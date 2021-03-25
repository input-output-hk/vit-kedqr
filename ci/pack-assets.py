# arguments: vit-kedqr version, target triple, target cpu

import sys
import platform
import hashlib
import shutil
import os


def sha256sum(path):
    h = hashlib.sha256()
    with open(path, 'rb') as f:
        data = f.read()
        h.update(data)
    return h.hexdigest()


version = sys.argv[1]
target = sys.argv[2]
target_cpu = sys.argv[3]

archive_basename = 'vit-kedqr-{}-{}-{}'.format(version, target, target_cpu)

root_dir = './target/{}/release'.format(target)

if platform.system() == 'Windows':
    vit_kedqr_name = 'vit-kedqr.exe'
else:
    vit_kedqr_name = 'vit-kedqr'

vit_kedqr_path = os.path.join(root_dir, vit_kedqr_name)

vit_kedqr_checksum = sha256sum(vit_kedqr_path)

# build archive
if platform.system() == 'Windows':
    import zipfile
    content_type = 'application/zip'
    archive_name = '{}.zip'.format(archive_basename)
    with zipfile.ZipFile(archive_name, mode='x') as archive:
        archive.write(vit_kedqr_path, arcname=vit_kedqr_name)
else:
    import tarfile
    content_type = 'application/gzip'
    archive_name = '{}.tar.gz'.format(archive_basename)
    with tarfile.open(archive_name, 'x:gz') as archive:
        archive.add(vit_kedqr_path, arcname=vit_kedqr_name)

# verify archive
shutil.unpack_archive(archive_name, './unpack-test')
vit_kedqr1_checksum = sha256sum(
    os.path.join('./unpack-test', vit_kedqr_name))
shutil.rmtree('./unpack-test')
if vit_kedqr1_checksum != vit_kedqr_checksum:
    print('vit_kedqr checksum mismarch: before {} != after {}'.format(
        vit_kedqr_checksum, vit_kedqr1_checksum))
    exit(1)

# save archive checksum
archive_checksum = sha256sum(archive_name)
checksum_filename = '{}.sha256'.format(archive_name)
with open(checksum_filename, 'x') as f:
    f.write(archive_checksum)

# set GitHub Action step outputs
print('::set-output name=release-archive::{}'.format(archive_name))
print('::set-output name=release-content-type::{}'.format(content_type))
