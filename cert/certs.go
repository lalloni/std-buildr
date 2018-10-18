package cert

const rootAFIP = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 5339851671891190767 (0x4a1af6cdc859a3ef)
    Signature Algorithm: sha512WithRSAEncryption
        Issuer: CN = AFIP Root CA 2, O = AFIP, C = AR
        Validity
            Not Before: Jan 18 15:28:52 2016 GMT
            Not After : Jan 18 03:00:00 2026 GMT
        Subject: CN = AFIP SSL CA
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (2048 bit)
                Modulus:
                    00:86:06:9d:be:99:2c:d8:fb:38:ea:c6:19:ee:42:
                    82:c6:53:04:7e:aa:50:f6:d3:37:9d:a2:fd:a0:39:
                    94:cc:b6:60:93:6f:c3:84:70:8d:8b:81:bf:c0:37:
                    7d:48:48:90:c0:dc:ca:f8:51:67:33:e6:c5:88:8d:
                    6d:76:c8:54:c5:38:1b:ab:4e:4f:e7:bd:b7:ee:a0:
                    6b:9f:65:ed:db:1a:c7:0b:56:2e:f3:13:71:a9:d7:
                    c8:b5:7b:b0:86:7a:46:cc:58:8d:97:d5:4a:3d:d3:
                    59:75:27:50:7e:54:aa:78:f9:3f:5a:c4:2b:1a:ad:
                    83:fc:21:6e:99:f6:59:89:cf:5f:71:92:11:37:bb:
                    52:fd:cc:2a:dc:39:53:a4:73:da:c7:34:60:5a:b1:
                    62:c7:7b:39:88:32:32:8b:8c:50:24:14:40:d9:d2:
                    d6:fc:04:26:19:31:ce:80:f6:9e:7b:6e:7f:89:cd:
                    e0:25:6b:08:a2:db:01:18:37:8b:b6:34:00:db:9c:
                    e6:bb:9f:c8:47:e3:bc:cb:f1:f5:1d:1a:09:b8:7f:
                    d6:d6:6e:da:92:a1:90:81:35:e8:6b:37:aa:2f:57:
                    46:98:b4:ff:c2:9a:c4:60:34:14:76:14:32:16:e0:
                    2b:d2:7d:b1:0a:16:64:24:29:d8:8a:c3:67:3b:45:
                    89:1f
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            Authority Information Access:
                OCSP - URI:http://auth.afip.gob.ar/ocsp/

            X509v3 Subject Key Identifier:
                38:2A:5E:8A:11:C9:7B:B1:44:27:61:4B:E4:1E:D1:51:59:6C:DC:87
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Authority Key Identifier:
                keyid:AE:6F:A6:EC:94:79:74:6B:50:34:DA:99:36:55:B3:71:5D:28:9A:06

            X509v3 CRL Distribution Points:

                Full Name:
                  URI:http://auth.afip.gob.ar/crl/afip-root-ca-2.crl

            X509v3 Key Usage: critical
                Digital Signature, Certificate Sign, CRL Sign
    Signature Algorithm: sha512WithRSAEncryption
         71:ed:c9:80:82:c4:75:19:f8:36:73:91:28:ef:b6:64:56:46:
         0e:f7:70:37:9a:b6:50:25:aa:15:df:7e:96:45:e4:67:a3:76:
         c4:71:25:f1:c1:df:14:23:e2:d2:1f:8c:62:ed:7f:9e:f9:86:
         aa:9a:54:b2:c1:05:75:ba:fe:10:5c:a7:30:5e:df:02:ef:96:
         86:19:26:16:ac:97:6c:3b:3b:66:a0:bc:3b:a4:ec:2c:5c:07:
         4b:d7:8e:03:73:78:f7:ba:93:a0:f6:cc:a2:2f:3d:c2:47:ad:
         de:86:7a:80:80:08:64:42:dc:15:39:3f:03:a0:c0:d8:17:60:
         a7:7a:fb:2c:cd:80:03:d9:61:e9:e4:a2:bf:61:fd:63:23:d3:
         bb:4f:a3:f9:98:6c:85:68:6d:33:5f:2c:52:cf:24:f9:b6:7e:
         80:bb:64:13:00:35:87:a4:a7:55:60:f0:4b:9e:35:9d:df:d5:
         d8:dd:41:6d:54:67:e8:33:97:0b:0d:c3:23:fd:94:66:9c:5e:
         fb:9d:5a:d0:a8:b8:2b:3c:1e:c0:10:eb:fd:c1:e4:91:94:04:
         30:07:1a:af:b2:9a:4f:2c:cc:25:ab:23:e4:13:d9:ae:62:a3:
         25:be:8e:fe:d3:ad:10:82:7a:eb:13:d2:05:3c:4f:4a:19:04:
         60:79:c1:ec
-----BEGIN CERTIFICATE-----
MIIDrjCCApagAwIBAgIIShr2zchZo+8wDQYJKoZIhvcNAQENBQAwNTEXMBUGA1UE
AwwOQUZJUCBSb290IENBIDIxDTALBgNVBAoMBEFGSVAxCzAJBgNVBAYTAkFSMB4X
DTE2MDExODE1Mjg1MloXDTI2MDExODAzMDAwMFowFjEUMBIGA1UEAwwLQUZJUCBT
U0wgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCGBp2+mSzY+zjq
xhnuQoLGUwR+qlD20zedov2gOZTMtmCTb8OEcI2Lgb/AN31ISJDA3Mr4UWcz5sWI
jW12yFTFOBurTk/nvbfuoGufZe3bGscLVi7zE3Gp18i1e7CGekbMWI2X1Uo901l1
J1B+VKp4+T9axCsarYP8IW6Z9lmJz19xkhE3u1L9zCrcOVOkc9rHNGBasWLHezmI
MjKLjFAkFEDZ0tb8BCYZMc6A9p57bn+JzeAlawii2wEYN4u2NADbnOa7n8hH47zL
8fUdGgm4f9bWbtqSoZCBNehrN6ovV0aYtP/CmsRgNBR2FDIW4CvSfbEKFmQkKdiK
w2c7RYkfAgMBAAGjgeAwgd0wOQYIKwYBBQUHAQEELTArMCkGCCsGAQUFBzABhh1o
dHRwOi8vYXV0aC5hZmlwLmdvYi5hci9vY3NwLzAdBgNVHQ4EFgQUOCpeihHJe7FE
J2FL5B7RUVls3IcwDwYDVR0TAQH/BAUwAwEB/zAfBgNVHSMEGDAWgBSub6bslHl0
a1A02pk2VbNxXSiaBjA/BgNVHR8EODA2MDSgMqAwhi5odHRwOi8vYXV0aC5hZmlw
LmdvYi5hci9jcmwvYWZpcC1yb290LWNhLTIuY3JsMA4GA1UdDwEB/wQEAwIBhjAN
BgkqhkiG9w0BAQ0FAAOCAQEAce3JgILEdRn4NnORKO+2ZFZGDvdwN5q2UCWqFd9+
lkXkZ6N2xHEl8cHfFCPi0h+MYu1/nvmGqppUssEFdbr+EFynMF7fAu+WhhkmFqyX
bDs7ZqC8O6TsLFwHS9eOA3N497qToPbMoi89wket3oZ6gIAIZELcFTk/A6DA2Bdg
p3r7LM2AA9lh6eSiv2H9YyPTu0+j+ZhshWhtM18sUs8k+bZ+gLtkEwA1h6SnVWDw
S541nd/V2N1BbVRn6DOXCw3DI/2UZpxe+51a0Ki4KzwewBDr/cHkkZQEMAcar7Ka
TyzMJasj5BPZrmKjJb6O/tOtEIJ66xPSBTxPShkEYHnB7A==
-----END CERTIFICATE-----
`
