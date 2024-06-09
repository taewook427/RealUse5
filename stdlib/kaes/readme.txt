stdlib5.kaes [win & linux] : 데이터/파일 암호화.

<python-st>
class allmode
    func encrypt(bytes pw, bytes kf, any data, int pmode) -> any
    # 파일/바이트 암호화. data는 파일 경로 또는 바이트. 암호화 파일 경로 혹은 암호화 바이트 반환.
    func decrypt(bytes pw, bytes kf, any data) -> any
    # 파일/바이트 복호화. data는 파일 경로 또는 바이트. 평문 파일 경로 혹은 평문 바이트 반환.
    func view(any data)
    # 암호화된 파일/바이트 정보 보기. 서명을 검증하고 내부 값을 업데이트.

    .hint str # 비밀번호/키파일 힌트 문자열.
    .msg str # 프로그램 지시용 추가정보 문자열.
    .signkey str[2] # ksign 서명 정보. (공개키, 비밀키)
    .proc float # 진행 정보. -1 : 시작 전, 0~1 : 진행 중, 2 : 종료 후.

class funcmode
    func encrypt(bytes akey)
    # 파일/바이트 암호화. akey는 48바이트.
    func decrypt(bytes akey)
    # 파일/바이트 복호화. akey는 48바이트.

    .before any # 읽기 대상 파일 경로 / 바이트.
    .after any # 쓰기 대상 파일 경로 / 바이트.
    .proc float # 진행 정보. -1 : 시작 전, 0~1 : 진행 중, 2 : 종료 후.

func genrand(int size) -> bytes
# size 크기의 암호학적으로 안전한 난수 바이트 반환.
func basickey() -> bytes
# gen5kaes 기본키파일 반환.

<python-hy>
class allmode
    func encrypt(bytes pw, bytes kf, any data, int pmode) -> any
    # 파일/바이트 암호화. data는 파일 경로 또는 바이트. 암호화 파일 경로 혹은 암호화 바이트 반환.
    func decrypt(bytes pw, bytes kf, any data) -> any
    # 파일/바이트 복호화. data는 파일 경로 또는 바이트. 평문 파일 경로 혹은 평문 바이트 반환.
    func view(any data)
    # 암호화된 파일/바이트 정보 보기. 서명을 검증하고 내부 값을 업데이트.
    func genrand(int size) -> bytes
    # size 크기의 암호학적으로 안전한 난수 바이트 반환.
    func basickey() -> bytes
    # gen5kaes 기본키파일 반환.

    .hint str # 비밀번호/키파일 힌트 문자열.
    .msg str # 프로그램 지시용 추가정보 문자열.
    .signkey str[2] # ksign 서명 정보. (공개키, 비밀키)
    .proc float # 진행 정보. -1 : 시작 전, 0~1 : 진행 중, 2 : 종료 후.

class funcmode
    func encrypt(bytes akey)
    # 파일/바이트 암호화. akey는 48바이트.
    func decrypt(bytes akey)
    # 파일/바이트 복호화. akey는 48바이트.
    func genrand(int size) -> bytes
    # size 크기의 암호학적으로 안전한 난수 바이트 반환.
    func basickey() -> bytes
    # gen5kaes 기본키파일 반환.

    .before any # 읽기 대상 파일 경로 / 바이트.
    .after any # 쓰기 대상 파일 경로 / 바이트.
    .proc float # 진행 정보. -1 : 시작 전, 0~1 : 진행 중, 2 : 종료 후.

<go>
struct SimIO
    func Open(interface v, bool isreader) -> error
    # 파일 경로 문자열 혹은 데이터 바이트열, 읽기모드 여부를 받아 내부 정보를 초기화.
    func Close() -> byte[]
    # 내부 파일 객체를 닫음. 파일 모드는 nil, 바이트 모드는 버퍼 데이터 반환.
    func Seek(int pos)
    # 위치 pos로 이동. 보통 읽기모드에서만 사용.
    func Read(int size) -> byte[]
    # size만큼 읽기 시도, 반환 바이트열의 크기는 size 이하.
    func Write(byte[] data)
    # data 쓰기 시도.

    .Buf byte[] # 내부 읽기쓰기 버퍼.
    .File OSFILE* # 내부 파일 입출력 버퍼.
    .Size int # 읽기모드에서 데이터 전체 크기.

struct Allmode
    func EnBin(byte[] pw, byte[] kf, byte[] data, int pmode) -> (byte[], error)
    # 이진 데이터 암호화. 암호화 결과물 반환.
    func EnFile(byte[] pw, byte[] kf, str path, int pmode) -> (str, error)
    # 파일 암호화. 암호화 파일 경로 반환.
    func ViewBin(byte[] data) -> error
    # 암호화 바이트 정보 보기. 서명을 검증하고 내부 값을 업데이트.
    func ViewFile(str path) -> error
    # 암호화 파일 정보 보기. 서명을 검증하고 내부 값을 업데이트.
    func DeBin(byte[] pw, byte[] kf, byte[] data) -> (byte[], error)
    # 이진 데이터 복호화. 복호화 결과물 반환.
    func DeFile(byte[] pw, byte[] kf, str path) -> (str, error)
    # 파일 복호화. 복호화 파일 경로 반환.

    .Hint str # 비밀번호/키파일 힌트 문자열.
    .Msg str # 프로그램 지시용 추가정보 문자열.
    .Signkey str[2] # ksign 서명 정보. (공개키, 비밀키)
    .Proc float # 진행 정보. -1 : 시작 전, 0~1 : 진행 중, 2 : 종료 후.

struct Funcmode
    func Encrypt(byte[] akey) -> error
    # 파일/바이트 암호화. akey는 48바이트.
    func Decrypt(byte[] akey) -> error
    # 파일/바이트 복호화. akey는 48바이트.

    .Before SimIO # 읽기 대상 파일 경로 / 바이트 설정된 SimIO 구조체.
    .After SimIO # 쓰기 대상 파일 경로 / 바이트 설정된 SimIO 구조체.
    .Proc float # 진행 정보. -1 : 시작 전, 0~1 : 진행 중, 2 : 종료 후.

!! 내부 값 초기화 !!
모든 작업 전후로 내부 값을 초기화해야
값 오염을 막을 수 있습니다.
특히 서명 정보를 초기화하십시오.

!!! SimIO 대용량 입출력 경고 !!!
SimIO는 go IO 모델의 얇은 추상화 버전입니다.
go file RW의 문제점인 1GiB 이상의 큰 파일
입출력이 제대로 되지 않는 문제가 있습니다.
또한 더 빠른 바이트 쓰기 작업을 위해,
Open시 미리 용량을 할당한 슬라이스를 사용하십시오.

!!! python VM 경고 !!!
threading 모듈 사용 시 메인 스크립트에
if __name__ == '__main__':
문구를 추가하세요.

!!! python VM 경고 !!!
multiprocessing 모듈 사용 시 메인 스크립트에
import multiprocessing as mp
if __name__ == '__main__':
    mp.freeze_support()
문구를 추가하세요.

전자서명 기능이 통합된 암호화 라이브러리입니다.
all-mode는 사진 헤더를 포함하는 사용자용 기능이며,
func-mode는 랜덤한 난수로 보이는 프로그램용 기능입니다.

all-mode에서 Bmode/Fmode 모두 원본 파일 이름이 포함되는데,
Bmode에서는 NewData.bin으로 고정됩니다.
또한 Fmode에서는 원본 파일과 같은 폴더 상에 암호화 파일이 생성되며,
이름은 16진법 표시된 2바이트 난수에 헤더에 맞는 확장자로 설정됩니다.

all-mode 암호화 시 입력받는 pmode는 사진 헤더를 결정합니다.
0 : webp, 1 : png, else : None 입니다.
각각 .webp .png .k 확장자를 가집니다.

go 버전의 func-mode의 구조체 내부 IO 정보는
함수 실행 전에 호출자가 미리 설정해야 합니다.
kaes 라이브러리는 IO 정보를 Open/Close 하지 않으니
입출력의 책임은 모두 호출자가 집니다.

KAES5 all-mode는 KSC5 청크를 사용합니다.
G5KAES all-mode {
    preheader : 512nB // webp/png/empty 사진헤더.
    common : 4B; KSC5 // KSC5 공통 표지.
    subtype : 4B; KAES // KAES5 표지.
    reserved : 8B; c0/c1 CRC32 // 청크0/1의 CRC32 값을 이어붙인 것.

    data chunk {
        c0 : encheader {
            // KDB5 형식의 텍스트를 utf-8 인코딩.
            salt; 40B
            pwhash; 128B
            ckeydata; 1920B
            tkeydata; 48B
            encname; 16nB
            hint; str
            msg; str
        }

        c1 : signheader {
            // 전자서명 기능을 사용하지 않을 경우 empty.
            publickey; str // 공개키 텍스트.
            signdata; bytes // c0 바이너리의 SHA3-512 해시값이 암호화된 바이너리.
        }

        c2 : encdata // 16nB 사이즈의 암호화된 데이터.
    }
    endsign : 8B; 8xFF
}

KAES5 func-mode는 겉보기에 무작위 바이트와 구분이 불가합니다.
G5KAES func-mode {
    ckeydata : 1920B // akey로 암호화된 ckey.
    encdata : 16nB // 암호화된 데이터.
}

비밀번호 흐름 (모든 키는 (iv16 + key32) * n 형식)
< all-mode >
pw, kf, salt(R40) -> pwh(128)
pw, kf, salt(R40) -> mkey(96) -> [ ckey(R1920), tkey(R48) ] -> ckeydata(1920), tkeydata(48)
ckey -> data, tkey -> name
< func-mode >
akey(48) -> [ ckey(R1920) ] -> ckeydata(1920)
ckey -> data
