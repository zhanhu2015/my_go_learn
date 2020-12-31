import ctypes
from ctypes import CDLL
_lib = CDLL('./libencrypt.so')
enc = _lib.EncryptFile
enc.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
if __name__ == '__main__':
    enc(b"/home/megvii/beeworker-nvr/tmp/test.jpg", b"/home/megvii/beeworker-nvr/tmp/test.jpg.meg")

