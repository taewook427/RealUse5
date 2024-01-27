stdlib5.kaes : 데이터/파일 암호화.

<python.st>
func genrandom(int size) -> bytes secure
# size 바이트의 안전한 난수 생성.
func genkf(str path) -> bytes keyfile
# path 경로의 키 파일 바이트 반환. 읽기 실패 시 기본키파일 반환.

class genbytes
    func en(bytes pw, bytes kf, bytes hint, bytes data) -> bytes enc
    # 일반 바이트열 암호화. (사용자용)
    func de(bytes pw, bytes kf, bytes data, int stpoint) -> bytes plain
    # 일반 바이트열 복호화. (사용자용)
    func view(bytes data) -> (bytes hint, str msg, int stpoint)
    # KSC5-KAES 바이트열을 해석.

    .valid bool # 이 모듈인 유효한지 여부. (True가 유효)
    .noerr bool # CRC오류를 무시할지 여부. (False시 오류 정지)
    .mode str # 출력물의 fake header 모드. "webp"/"png"/""
    .msg str # 프로그램용 추가 메세지 설정.

class genfile
    func en(bytes pw, bytes kf, bytes hint, str path) -> str result
    # 일반 파일 암호화 (사용자용)
    func de(bytes pw, bytes kf, str path, int stpoint) -> str result
    # 일반 파일 복호화 (사용자용)
    func view(str path) -> (bytes hint, str msg, int stpoint)
    # KSC5-KAES 파일을 해석.

    .valid bool # 이 모듈인 유효한지 여부. (True가 유효)
    .noerr bool # CRC오류를 무시할지 여부. (False시 오류 정지)
    .mode str # 출력물의 fake header 모드. "webp"/"png"/""
    .msg str # 프로그램용 추가 메세지 설정.

class funcbytes
    func en(bytes key, bytes data) -> bytes enc
    # 바이트 열 암호화, 키는 48B. (프로그램 내부 이용)
    func de(bytes key, bytes data) -> bytes plain
    # 바이트 열 복호화, 키는 48B. (프로그램 내부 이용)

    .valid bool # 이 모듈인 유효한지 여부. (True가 유효)

class funcfile
    func en(bytes key, str before, str after)
    # 파일 암호화, 키는 48B. (프로그램 내부 이용)
    func de(bytes key, str before, str after)
    # 파일 복호화, 키는 48B. (프로그램 내부 이용)

    .valid bool # 이 모듈인 유효한지 여부. (True가 유효)

<python.hy>
standard 버전과 동일.

<go>
func Genrandom(int size) -> byte[] secure
# size 바이트의 안전한 난수 생성.
func Genkf(string path) -> bytes[] keyfile
# path 경로의 키 파일 바이트 반환. 읽기 실패 시 기본키파일 반환.

func Init0() -> genbytes
# 일반 사용자용 바이트열 모듈 생성.
struct genbytes
    func En(byte[] pw, byte[] kf, byte[] hint, byte[] data) -> byte[] enc
    # 일반 바이트열 암호화. (사용자용)
    func De(byte[] pw, byte[] kf, byte[] data, int stpoint) -> byte[] plain
    # 일반 바이트열 복호화. (사용자용)
    func View(byte[] data) -> (byte[] hint, string msg, int stpoint)
    # KSC5-KAES 바이트열을 해석.

    .Valid bool # 이 모듈인 유효한지 여부. (True가 유효)
    .Noerr bool # CRC오류를 무시할지 여부. (False시 오류 정지)
    .Mode string # 출력물의 fake header 모드. "webp"/"png"/""
    .Msg string # 프로그램용 추가 메세지 설정.

func Init1() -> genfile
# 일반 사용자용 파일 모듈 생성.
struct genfile
    func En(byte[] pw, byte[] kf, byte[] hint, string path) -> string result
    # 일반 파일 암호화 (사용자용)
    func De(byte[] pw, byte[] kf, string path, int stpoint) -> string result
    # 일반 파일 복호화 (사용자용)
    func View(string path) -> (byte[] hint, string msg, int stpoint)
    # KSC5-KAES 파일을 해석.

    .Valid bool # 이 모듈인 유효한지 여부. (True가 유효)
    .Noerr bool # CRC오류를 무시할지 여부. (False시 오류 정지)
    .Mode string # 출력물의 fake header 모드. "webp"/"png"/""
    .Msg string # 프로그램용 추가 메세지 설정.

func Init2() -> funcbytes
# 기능적 바이트열 모듈 생성.
struct funcbytes
    func En(byte[] key, byte[] data) -> byte[] enc
    # 바이트 열 암호화, 키는 48B. (프로그램 내부 이용)
    func De(byte[] key, byte[] data) -> byte[] plain
    # 바이트 열 복호화, 키는 48B. (프로그램 내부 이용)

    .Valid bool # 이 모듈인 유효한지 여부. (True가 유효)

func Init3() -> funcfile
# 기능적 파일 모듈 생성.
struct funcfile
    func En(byte[] key, string before, string after)
    # 파일 암호화, 키는 48B. (프로그램 내부 이용)
    func De(byte[] key, string before, string after)
    # 파일 복호화, 키는 48B. (프로그램 내부 이용)

    .Valid bool # 이 모듈인 유효한지 여부. (True가 유효)

!!! python VM 경고 !!!
multiprocessing 모듈 사용 시 메인 스크립트에
import multiprocessing as mp
if __name__ == '__main__':
    mp.freeze_support()
문구를 추가할 것.

kaes5는 사용자용 general 버전과 subsystem용 functional 버전이 있습니다.
gen은 KAES4와 비슷한 구조로 이루어져 있으며, 분할 암호화로 32코어까지 병렬 연산을 지원합니다.
func는 필수적인 암호화 기능만 남겨놓은 버전이라 키 없이는 랜덤  바이트열로 보입니다.

KSC5 {
    prehead + pad : 1024nB // 사진 정보를 담을 수 있는 가짜 헤더, 0또는 1KiB 배수 크기.
    common magicnum : 4B // KSC5, KSC5 파일임을 나타냄.
    subtype magicnum : 4B // KSC5 파일 중 어느 종류인지 나타냄.
    reserved : 4B // 다양한 방법으로 사용할 수 있게 예약된 영역.
    mainheader size : 4B // 리틀 엔디안 인코딩된 메인헤더의 크기.
    mainheader : nB // 메인 헤더.
    (
        chunk size : 8B // 리틀 엔디안 인코딩된 청크의 크기.
        chunk : nB // 청크. 청크는 0개 혹은 여러 개가 올 수 있다.
    )...
    trash : nB // 청크가 모두 끝나고 남은 부분이 있다면 쓰레기값이다.
}

KAES gen은 KSC5의 서브타입으로 구성됩니다.
서브타입은 KAES, 예약 영역은 mainheader의 CRC32 체크섬입니다.
청크는 단일청크만 사용합니다.
mainheader는 UTF-8인코딩된 KDB 문자열입니다.

KAES gen MH {
    mode : str // "bytes"/"file"로 암호화 종류를 표시.
    msg : str // 추가 메세지 전달 문자열. (프로그램용)
    salt : 32B // 비밀번호 해시 보안용 salt.
    pwhash : 256B // 비밀번호 저장용 해시.
    hint : nB // 비밀번호 힌트. (사용자용)
    tkeydt : 48B // 파일 이름 암호화 키. (file mode only)
    ckeydt : 1536B // 32 병렬 데이터 암호화 키.
    namedt : 16nB // 패딩된 원본 이름 암호문. (file mode only)
}

비밀번호 흐름 : 32B rand SALT. PW + KF + SALT -> masterkey MKEY 48B, pwsave PWHASH 256B.
48B rand titlekey TKEY. TKEY -(MKEY)-> tkey data TKEYDT 48B. (file mode only)
nB NAME -(padding)-> 16nB NMB -(TKEY)-> 16nB name data NAMEDT. (file mode only)
1536B rand contentkey CKEY. CKEY -(MKEY)-> ckey data CKEYDT 1536B.
user data -(CKEY)-> encrypted data.
데이터 암호화 방식 : 32스레드 병렬. 청크 크기는 512KiB 고정. 16MiB부터 32스레드 전체 이용 시작.
데이터를 512KiB 단위로 나눠서 처음 32개 청크를 동시에, 그 후 다음 32개 청크를 동시에... 진행함.
즉, 독립적인 32개의 암호화 흐름을 만들어 암호화. 16바이트 패딩은 원본 데이터의 끝인 마지막 청크 1개에서만 진행함.
마지막 청크의 크기는 512KiB보다 작거나 같고, 패딩된 마지막 청크의 길이는 최소 16바이트.

KAES func {
    ckey data : 1536B // 32 분할 암호화에서 사용하는 키와 초기화 벡터.
}

<벤치마크>
환경 : os = windows 11, cpu = i5-13600k, ram = 32 GiB
측정값은 MiB/s 단위로, 7회 반복하여 얻은 결과의 평균/최대/최소.

===== 평균 =====
MiB/s    genBen    genBde    genFen    genFde   funcBen   funcBde   funcFen   funcFde   randgen

py.st      95.6      95.2     103.9     104.2      97.7      98.0     107.2     105.9    3624.7

py.hy      94.2      93.6    1354.5    1140.7      95.6      95.0    1641.9    1321.0    3764.1

go.st     868.3    1455.6    1397.2    1260.7    2118.2    3829.9    1674.1    1355.4    3480.8

===== 최대 =====
MiB/s    genBen    genBde    genFen    genFde   funcBen   funcBde   funcFen   funcFde   randgen

py.st      97.5      96.1     105.3     105.2      98.5      99.4     108.1     107.0    3797.1

py.hy      95.5      94.3    1417.8    1222.4      97.0      96.6    1891.2    1452.5    3855.1

go.st     983.8    1627.6    1480.0    1306.6    2265.8    4169.4    1923.8    1724.0    3964.1

===== 최소 =====
MiB/s    genBen    genBde    genFen    genFde   funcBen   funcBde   funcFen   funcFde   randgen

py.st      93.7      94.0     102.0     103.7      96.6      95.4     106.3     100.7    3472.0

py.hy      92.8      93.1    1306.3     937.2      94.5      94.0    1441.8    1179.7    3618.2

go.st     617.4    1268.6    1269.5    1164.1    1956.9    3166.4    1422.2    1173.7    1970.3

<원본 데이터>
===== gen5 KAES py.st test =====
k0 valid
k1 valid
k2 valid
k3 valid

test 0 : hint_True, msg_True
test 0 : data_True
test 1 : hint_True, msg_True
test 1 : data_True
test 2 : hint_True, msg_True
test 2 : data_True
test 3 : hint_True, msg_True
test 3 : data_True
test 4 : hint_True, msg_True
test 4 : data_True
test 5 : hint_True, msg_True
test 5 : data_True
test 6 : hint_True, msg_True
test 6 : data_True
test 7 : hint_True, msg_True
test 7 : data_True

test 0 : hint_True, msg_True
test 0 : data_True
test 1 : hint_True, msg_True
test 1 : data_True
test 2 : hint_True, msg_True
test 2 : data_True
test 3 : hint_True, msg_True
test 3 : data_True
test 4 : hint_True, msg_True
test 4 : data_True
test 5 : hint_True, msg_True
test 5 : data_True
test 6 : hint_True, msg_True
test 6 : data_True
test 7 : hint_True, msg_True
test 7 : data_True

test 0 : data_True
test 1 : data_True
test 2 : data_True
test 3 : data_True
test 4 : data_True
test 5 : data_True
test 6 : data_True
test 7 : data_True

test 0 : data_True
test 1 : data_True
test 2 : data_True
test 3 : data_True
test 4 : data_True
test 5 : data_True
test 6 : data_True
test 7 : data_True

rand : time 0.5614399909973145 s, speed 3647.7629539036457 MiB/s
k0en : time 21.19334363937378 s, speed 96.63411469416037 MiB/s
k0de : time 21.51251530647278 s, speed 95.20039710948114 MiB/s
k1en : time 39.24112391471863 s, speed 104.38029269757142 MiB/s
k1de : time 39.16727256774902 s, speed 104.57710561579194 MiB/s
k2en : time 10.412832021713257 s, speed 98.34020157673858 MiB/s
k2de : time 10.390108823776245 s, speed 98.55527188095718 MiB/s
k3en : time 38.516422510147095 s, speed 106.3442483247481 MiB/s
k3de : time 40.66514325141907 s, speed 100.72508474089943 MiB/s

rand : time 0.5632641315460205 s, speed 3635.9496110266196 MiB/s
k0en : time 21.012338638305664 s, speed 97.46654264682749 MiB/s
k0de : time 21.309123039245605 s, speed 96.1090701024224 MiB/s
k1en : time 40.16230630874634 s, speed 101.98617500977514 MiB/s
k1de : time 38.92041206359863 s, speed 105.2404068411931 MiB/s
k2en : time 10.46785593032837 s, speed 97.82327984025645 MiB/s
k2de : time 10.433001279830933 s, speed 98.15008860198223 MiB/s
k3en : time 37.8869743347168 s, speed 108.11103477974834 MiB/s
k3de : time 38.32753849029541 s, speed 106.8683291789561 MiB/s

rand : time 0.5681822299957275 s, speed 3604.4773875018936 MiB/s
k0en : time 21.56445336341858 s, speed 94.97110663951158 MiB/s
k0de : time 21.384857892990112 s, speed 95.76869812500965 MiB/s
k1en : time 38.88507533073425 s, speed 105.33604384617394 MiB/s
k1de : time 39.18942904472351 s, speed 104.51798099241479 MiB/s
k2en : time 10.468528509140015 s, speed 97.81699492015055 MiB/s
k2de : time 10.395734310150146 s, speed 98.50194026218917 MiB/s
k3en : time 38.125848054885864 s, speed 107.43367581236252 MiB/s
k3de : time 38.296297550201416 s, speed 106.95550907057482 MiB/s

rand : time 0.5898630619049072 s, speed 3471.9922847621224 MiB/s
k0en : time 21.460081100463867 s, speed 95.43300374366859 MiB/s
k0de : time 21.56429648399353 s, speed 94.97179755064874 MiB/s
k1en : time 39.64525270462036 s, speed 103.31627926595715 MiB/s
k1de : time 39.37543964385986 s, speed 104.0242353367278 MiB/s
k2en : time 10.594965696334839 s, speed 96.64967583181856 MiB/s
k2de : time 10.561395168304443 s, speed 96.9568871992502 MiB/s
k3en : time 38.28812527656555 s, speed 106.97833781135736 MiB/s
k3de : time 38.36829328536987 s, speed 106.75481365656253 MiB/s

rand : time 0.5756187438964844 s, speed 3557.910547069154 MiB/s
k0en : time 21.857821226119995 s, speed 93.69643839673506 MiB/s
k0de : time 21.794008016586304 s, speed 93.97078309053443 MiB/s
k1en : time 39.23970603942871 s, speed 104.38406434248694 MiB/s
k1de : time 39.472453355789185 s, speed 103.76856900887982 MiB/s
k2en : time 10.395658254623413 s, speed 98.50266091082607 MiB/s
k2de : time 10.29819631576538 s, speed 99.4348882660521 MiB/s
k3en : time 38.02514624595642 s, speed 107.71819189086135 MiB/s
k3de : time 38.270904779434204 s, speed 107.02647412195712 MiB/s

rand : time 0.559962272644043 s, speed 3657.3892564756293 MiB/s
k0en : time 21.527273178100586 s, speed 95.13513314279876 MiB/s
k0de : time 21.38458561897278 s, speed 95.7699174765855 MiB/s
k1en : time 39.33635878562927 s, speed 104.12758390581868 MiB/s
k1de : time 39.470136880874634 s, speed 103.77465911410934 MiB/s
k2en : time 10.528157949447632 s, speed 97.26297847324042 MiB/s
k2de : time 10.311852931976318 s, speed 99.30320057461732 MiB/s
k3en : time 38.45787167549133 s, speed 106.50615391725705 MiB/s
k3de : time 38.417091608047485 s, speed 106.61921109983201 MiB/s

rand : time 0.5393595695495605 s, speed 3797.095881158393 MiB/s
k0en : time 21.401968717575073 s, speed 95.69213127193312 MiB/s
k0de : time 21.666745901107788 s, speed 94.52273125588688 MiB/s
k1en : time 39.457284450531006 s, speed 103.80846165770227 MiB/s
k1de : time 39.48609662055969 s, speed 103.73271481758182 MiB/s
k2en : time 10.491642713546753 s, speed 97.60149368009041 MiB/s
k2de : time 10.739272117614746 s, speed 95.3509687421382 MiB/s
k3en : time 38.18541193008423 s, speed 107.26609437917264 MiB/s
k3de : time 38.47077775001526 s, speed 106.47042351511533 MiB/s

===== gen5 KAES py.hy test =====
k0 valid
k1 valid
k2 valid
k3 valid

test 0 : hint_True, msg_True
test 0 : data_True
test 1 : hint_True, msg_True
test 1 : data_True
test 2 : hint_True, msg_True
test 2 : data_True
test 3 : hint_True, msg_True
test 3 : data_True
test 4 : hint_True, msg_True
test 4 : data_True
test 5 : hint_True, msg_True
test 5 : data_True
test 6 : hint_True, msg_True
test 6 : data_True
test 7 : hint_True, msg_True
test 7 : data_True

test 0 : hint_True, msg_True
test 0 : data_True
test 1 : hint_True, msg_True
test 1 : data_True
test 2 : hint_True, msg_True
test 2 : data_True
test 3 : hint_True, msg_True
test 3 : data_True
test 4 : hint_True, msg_True
test 4 : data_True
test 5 : hint_True, msg_True
test 5 : data_True
test 6 : hint_True, msg_True
test 6 : data_True
test 7 : hint_True, msg_True
test 7 : data_True

test 0 : data_True
test 1 : data_True
test 2 : data_True
test 3 : data_True
test 4 : data_True
test 5 : data_True
test 6 : data_True
test 7 : data_True

test 0 : data_True
test 1 : data_True
test 2 : data_True
test 3 : data_True
test 4 : data_True
test 5 : data_True
test 6 : data_True
test 7 : data_True

rand : time 0.5312454700469971 s, speed 3855.091695782031 MiB/s
k0en : time 21.66318655014038 s, speed 94.53826173078534 MiB/s
k0de : time 21.72426176071167 s, speed 94.2724785108145 MiB/s
k1en : time 3.1356148719787598 s, speed 1306.2828718551077 MiB/s
k1de : time 3.472478151321411 s, speed 1179.561057408328 MiB/s
k2en : time 10.650030136108398 s, speed 96.14996266800962 MiB/s
k2de : time 10.825011253356934 s, speed 94.59574461711976 MiB/s
k3en : time 2.165863275527954 s, speed 1891.1627738835698 MiB/s
k3de : time 3.2715959548950195 s, speed 1251.988343447941 MiB/s

rand : time 0.5660271644592285 s, speed 3618.200907295006 MiB/s
k0en : time 21.450679540634155 s, speed 95.47483081458846 MiB/s
k0de : time 21.75627303123474 s, speed 94.13376992740237 MiB/s
k1en : time 3.1320059299468994 s, speed 1307.7880730798759 MiB/s
k1de : time 3.3598878383636475 s, speed 1219.0883139702837 MiB/s
k2en : time 10.553265571594238 s, speed 97.0315769136196 MiB/s
k2de : time 10.832632303237915 s, speed 94.52919395167899 MiB/s
k3en : time 2.840837001800537 s, speed 1441.8285869284066 MiB/s
k3de : time 2.8232743740081787 s, speed 1450.7977112351796 MiB/s

rand : time 0.5353517532348633 s, speed 3825.522168602155 MiB/s
k0en : time 21.523327350616455 s, speed 95.15257407174745 MiB/s
k0de : time 22.00832962989807 s, speed 93.05567639344218 MiB/s
k1en : time 2.9224343299865723 s, speed 1401.5712715839948 MiB/s
k1de : time 3.4562203884124756 s, speed 1185.1096109879122 MiB/s
k2en : time 10.670617580413818 s, speed 95.9644549420998 MiB/s
k2de : time 10.8888099193573 s, speed 94.04149834405783 MiB/s
k3en : time 2.520225763320923 s, speed 1625.251221383701 MiB/s
k3de : time 2.862388849258423 s, speed 1430.9725951668574 MiB/s

rand : time 0.5475413799285889 s, speed 3740.356574085968 MiB/s
k0en : time 21.785630702972412 s, speed 94.00691804256888 MiB/s
k0de : time 21.98597526550293 s, speed 93.15029127743139 MiB/s
k1en : time 3.1104722023010254 s, speed 1316.8418598854262 MiB/s
k1de : time 3.40278959274292 s, speed 1203.7182694855655 MiB/s
k2en : time 10.700790166854858 s, speed 95.69386783901125 MiB/s
k2de : time 10.604188680648804 s, speed 96.56561485638785 MiB/s
k3en : time 2.523261070251465 s, speed 1623.2961576155092 MiB/s
k3de : time 3.3067798614501953 s, speed 1238.667274997765 MiB/s

rand : time 0.5473318099975586 s, speed 3741.7887332532987 MiB/s
k0en : time 21.900712966918945 s, speed 93.51293736845493 MiB/s
k0de : time 21.876662254333496 s, speed 93.6157433977076 MiB/s
k1en : time 3.0208473205566406 s, speed 1355.9109631681897 MiB/s
k1de : time 3.3506569862365723 s, speed 1222.446826644762 MiB/s
k2en : time 10.769240140914917 s, speed 95.08563153955303 MiB/s
k2de : time 10.738187074661255 s, speed 95.36060350599759 MiB/s
k3en : time 2.4462039470672607 s, speed 1674.4311139349888 MiB/s
k3de : time 3.472130537033081 s, speed 1179.6791498225216 MiB/s

rand : time 0.534691572189331 s, speed 3830.2455219451554 MiB/s
k0en : time 21.859068870544434 s, speed 93.69109050933656 MiB/s
k0de : time 21.869245767593384 s, speed 93.64749117387478 MiB/s
k1en : time 2.9776885509490967 s, speed 1375.5636057687957 MiB/s
k1de : time 4.3705291748046875 s, speed 937.1862848125352 MiB/s
k2en : time 10.83647608757019 s, speed 94.49566369408254 MiB/s
k2de : time 10.71796441078186 s, speed 95.5405299694684 MiB/s
k3en : time 2.537245035171509 s, speed 1614.3493999282275 MiB/s
k3de : time 3.297215223312378 s, speed 1242.2604296619631 MiB/s

rand : time 0.5480039119720459 s, speed 3737.199598867955 MiB/s
k0en : time 22.078293561935425 s, speed 92.76079214432134 MiB/s
k0de : time 21.996899604797363 s, speed 93.1040299676299 MiB/s
k1en : time 2.8890466690063477 s, speed 1417.7687207139402 MiB/s
k1de : time 3.9458870887756348 s, speed 1038.0428805607166 MiB/s
k2en : time 10.766824007034302 s, speed 95.10696927255326 MiB/s
k2de : time 10.821630716323853 s, speed 94.62529510042795 MiB/s
k3en : time 2.5237784385681152 s, speed 1622.963385931729 MiB/s
k3de : time 2.8200595378875732 s, speed 1452.451604290666 MiB/s

===== gen5 KAES go.st test =====

test 0 : hint_true, msg_true
test 0 : data_true
test 1 : hint_true, msg_true
test 1 : data_true
test 2 : hint_true, msg_true
test 2 : data_true
test 3 : hint_true, msg_true
test 3 : data_true
test 4 : hint_true, msg_true
test 4 : data_true
test 5 : hint_true, msg_true
test 5 : data_true
test 6 : hint_true, msg_true
test 6 : data_true
test 7 : hint_true, msg_true
test 7 : data_true

test 0 : hint_true, msg_true
test 0 : data_true
test 1 : hint_true, msg_true
test 1 : data_true
test 2 : hint_true, msg_true
test 2 : data_true
test 3 : hint_true, msg_true
test 3 : data_true
test 4 : hint_true, msg_true
test 4 : data_true
test 5 : hint_true, msg_true
test 5 : data_true
test 6 : hint_true, msg_true
test 6 : data_true
test 7 : hint_true, msg_true
test 7 : data_true

test 0 : data_true
test 1 : data_true
test 2 : data_true
test 3 : data_true
test 4 : data_true
test 5 : data_true
test 6 : data_true
test 7 : data_true

test 0 : data_true
test 1 : data_true
test 2 : data_true
test 3 : data_true
test 4 : data_true
test 5 : data_true
test 6 : data_true
test 7 : data_true

rand : time 1.039460 s, speed 1970.253786 MiB/s
k0en : time 3.317088 s, speed 617.409005 MiB/s
k0de : time 1.614429 s, speed 1268.559968 MiB/s
k1en : time 2.767613 s, speed 1479.975705 MiB/s
k1de : time 3.154098 s, speed 1298.628007 MiB/s
k2en : time 0.462465 s, speed 2214.221617 MiB/s
k2de : time 0.245599 s, speed 4169.398084 MiB/s
k3en : time 2.129146 s, speed 1923.776012 MiB/s
k3de : time 2.375832 s, speed 1724.027625 MiB/s

rand : time 0.560454 s, speed 3654.180361 MiB/s
k0en : time 2.081781 s, speed 983.773029 MiB/s
k0de : time 1.458939 s, speed 1403.759856 MiB/s
k1en : time 2.848843 s, speed 1437.776669 MiB/s
k1de : time 3.134805 s, speed 1306.620348 MiB/s
k2en : time 0.522001 s, speed 1961.682066 MiB/s
k2de : time 0.264166 s, speed 3876.350477 MiB/s
k3en : time 2.392427 s, speed 1712.068958 MiB/s
k3de : time 2.720814 s, speed 1505.431830 MiB/s

rand : time 0.523430 s, speed 3912.653077 MiB/s
k0en : time 2.330087 s, speed 878.937138 MiB/s
k0de : time 1.300838 s, speed 1574.369752 MiB/s
k1en : time 2.914129 s, speed 1405.565780 MiB/s
k1de : time 3.167874 s, speed 1292.980718 MiB/s
k2en : time 0.453147 s, speed 2259.752354 MiB/s
k2de : time 0.253951 s, speed 4032.273943 MiB/s
k3en : time 2.363283 s, speed 1733.182188 MiB/s
k3de : time 3.438663 s, speed 1191.160634 MiB/s

rand : time 0.571623 s, speed 3582.780959 MiB/s
k0en : time 2.234770 s, speed 916.425404 MiB/s
k0de : time 1.400471 s, speed 1462.365161 MiB/s
k1en : time 2.903338 s, speed 1410.789925 MiB/s
k1de : time 3.518705 s, speed 1164.064620 MiB/s
k2en : time 0.506474 s, speed 2021.821456 MiB/s
k2de : time 0.254291 s, speed 4026.882587 MiB/s
k3en : time 2.266184 s, speed 1807.443703 MiB/s
k3de : time 2.671795 s, speed 1533.051750 MiB/s

rand : time 0.534344 s, speed 3832.736963 MiB/s
k0en : time 2.332646 s, speed 877.972911 MiB/s
k0de : time 1.258332 s, speed 1627.551393 MiB/s
k1en : time 2.941557 s, speed 1392.459844 MiB/s
k1de : time 3.261681 s, speed 1255.794175 MiB/s
k2en : time 0.451941 s, speed 2265.782480 MiB/s
k2de : time 0.269460 s, speed 3800.192979 MiB/s
k3en : time 2.880123 s, speed 1422.161484 MiB/s
k3de : time 3.477219 s, speed 1177.952841 MiB/s

rand : time 0.593793 s, speed 3449.013377 MiB/s
k0en : time 2.200285 s, speed 930.788511 MiB/s
k0de : time 1.381392 s, speed 1482.562517 MiB/s
k1en : time 2.958858 s, speed 1384.317869 MiB/s
k1de : time 3.182927 s, speed 1286.865831 MiB/s
k2en : time 0.476849 s, speed 2147.430319 MiB/s
k2de : time 0.323397 s, speed 3166.386825 MiB/s
k3en : time 2.734346 s, speed 1497.981601 MiB/s
k3de : time 3.463716 s, speed 1182.544989 MiB/s

rand : time 0.516636 s, speed 3964.106257 MiB/s
k0en : time 2.346190 s, speed 872.904581 MiB/s
k0de : time 1.494553 s, speed 1370.309383 MiB/s
k1en : time 3.226511 s, speed 1269.482732 MiB/s
k1de : time 3.356779 s, speed 1220.217357 MiB/s
k2en : time 0.523280 s, speed 1956.887326 MiB/s
k2de : time 0.273960 s, speed 3737.771938 MiB/s
k3en : time 2.524629 s, speed 1622.416601 MiB/s
k3de : time 3.489884 s, speed 1173.677979 MiB/s
