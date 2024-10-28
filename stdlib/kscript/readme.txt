stdlib5.kscript [win & linux] : 간략화 스크립트 언어.

<python-runtime>
func tostr(any v) -> str
# 6-type 변수를 문자열로 반환.
func decode(bytes data) -> int
# 리틀엔디안 signed int 16/32/64 디코딩.
func testio(int mode, any[] args) -> any
# 테스트 기능 구현. 인터럽트 코드와 인자를 받아 동작하고 결과 반환.

class kvm
    func view(str path) -> (str, int, str)
    # 실행파일 경로를 받아 내부 값을 업데이트하고 (프로그램 정보, ABI 정보, 서명 공개키) 반환.
    func load(bool sign)
    # view한 파일을 로드해 스택과 명령어 초기화. sign이 True여야 디지털 서명 검증함.
    func run() -> int
    # 로드한 명령어를 실행하고 인터럽트 코드 반환.

    .callmem any[] # 외부호출 시 인자 저장 필드.
    .errmsg str # 예외발생 시 기록 필드.
    .maxstk int # 최대 스택 크기. (기본값 1048576)
    (.ma any # 외부호출의 결과값을 설정하는 필드)

<go-runtime>
struct Vunit
    func Set(interface v)
    # 6-type 값을 받아 내부 필드를 설정. 이외의 타입은 None으로 설정.
    func ToString() -> str
    # 변수값의 문자열 형태 반환. (float은 32비트 정밀도로 표현됨)

    .Vtype byte # 6-type 변수체의 타입 표시. (0:none, 1:bool, 2:int, 3:float, 4:string, 5:bytes)
    .Vbool bool # bool 값 저장 필드.
    .Vint int # int 값 저장 필드.
    .Vfloat float # float 값 저장 필드.
    .Vstring str # string 값 저장 필드.
    .Vbytes byte[] # bytes 값 저장 필드.

struct KVM
    func Init()
    # 구조체 초기화. View 실행 전 항상 초기화가 필요.
    func View(str path) -> (str, int, str, error)
    # 실행파일 경로를 받아 내부 값을 업데이트하고 (프로그램 정보, ABI 정보, 서명 공개키) 반환.
    func Load(bool sign) -> error
    # View한 파일을 로드해 스택과 명령어 초기화. sign이 True여야 디지털 서명 검증함.
    func Run() -> int
    # 로드한 명령어를 실행하고 인터럽트 코드 반환.
    func SetRet(Vunit* ma)
    # 외부호출의 결과값을 레지스터 MA에 설정.

    .CallMem Vunit[] # 외부호출 시 인자 저장 필드.
    .SafeMem bool # 안전한 메모리 접근 옵션. (기본값 True)
    .RunOne bool # 한 사이클만 실행 옵션. (기본값 False)
    .ErrHlt bool # 일반예외 발생 시 정지 옵션. (기본값 True)
    .ErrMsg str # 예외발생 시 기록 필드.
    .MaxStk int # 최대 스택 크기. (기본값 16777216)

func TestIO(int mode, Vunit[] v) -> Vunit*
# 테스트 기능 구현. 인터럽트 코드와 인자를 받아 동작하고 결과 반환.

<go-compiler>
struct Parser
    func Init()
    # 기본 문법으로 초기화.
    func Split(str code) -> error
    # 소스코드를 받아 토큰 단위로 분리해 내부 필드에 저장.
    func Parse() -> error
    # 내부 필드의 토큰들을 분류하고 조작함.
    func Structify() -> (Token[], error)
    # 내부의 괄호를 묶어 부분 구조화한 토큰 반환.
    func GenAST(Token[] tokens) -> (AST[], AST, error)
    # 토큰을 구조화된 트리로 변환하여 (내부함수들, 메인흐름) 반환.

    .Type_Function str[] # 외부 함수 이름들. (테스트 기능이 포함되므로 조작 필요)

struct Compiler
    func Init()
    # 기본 설정으로 초기화. 외부 함수 관련 필드에 테스트 기능을 포함시킴.
    func Compile(AST* mainflow, AST[] functions) -> (str, error)
    # 구문 트리를 받아 텍스트 형식 어셈블리 코드로 변환.

    .OptConst bool # 상수 접기 옵션. (기본값 True)
    .OptAsm bool # 단축 명령어 옵션. (기본값 True)
    .OuterNum int[str] # 외부 함수 이름과 인터럽트 코드 매칭.
    .OuterParms int[str] # 외부 함수 이름과 인자 개수 매칭.

struct Assembler
    func SetKey(str iconpath, str public, str private)
    # 아이콘 사진, 서명키 설정. (아이콘 경로가 빈 문자열이면 설정없이 지나감)
    func GenExe(str asm) -> (byte[], error)
    # 어셈블리 코드를 받아 KELF 형식 실행파일 생성.

    .Icon byte[] # 프로그램 아이콘 사진.
    .Info str # 프로그램 정보.
    .ABIf int # ABI 정보. (사용하는 외부호출 종류)

struct Token
    func Write(int indent) -> str
    # 토큰(하나의 단어/표현식)의 내부 값을 디버그 모드로 출력.

struct AST
    func Write(int indent) -> str
    # 구문트리(트리 형식인 코드의 구조)의 내부 값을 디버그 모드로 출력.

struct CLiteral
    .Ivalue interface # 컴파일 시 코드의 리터럴 부분의 값.

struct CTree
    func Set(AST ast, bool opt) # AST를 받아 컴파일 트리를 설정.

<kscript-testf>
func test.input(any q) -> str
# 질문 q를 출력하고 사용자 입력을 반환.
func test.print(any v)
# v의 문자열 형태를 출력. (줄바꿈 없음)
func test.read(str path, int size) -> bytes
# path 경로의 파일을 처음부터 size 만큼 읽어옴. (음수라면 전체읽기)
func test.write(str path, str|bytes data)
# path 경로의 파일을 만들고 data 기록.
func test.time() -> float
# 현재 유닉스 시간을 실수 형태로 반환.
func test.sleep(int|float t)
# t초만큼 실행을 멈춤.

간단한 매크로 제어용 스크립트 언어에서 출발한 kscript는
python, golang, c에서 영향을 받아 비슷한 문법을 가진
절차지향 동적 타이핑 언어가 되었습니다.
다른 고급 언어 위에서 가상머신인 KVM을 구현하여 실행됩니다.

!! 런타임에 따라 문자열 출력 결과나 run 사이클 조건이 달라질 수 있습니다. !!
(python 런타임은 safemem, runone, errhlt 옵션이 모두 True로 고정되었습니다)

! kvm 자체에서 오류가 나더라도, 런타임 에러 없이 c_err 상태를 반환합니다. !

현재 python과 go의 런타임만 배포되지만 직접 런타임을 구축할 경우
code.txt를 컴파일하고 실행한 결과를 stdout.txt와 비교하십시오.
kscript-kvm의 자세한 문법이나 특성은 기술 문서를 참조하십시오.

KELF format은 KSC5 청크를 사용합니다.
KELF5 {
    preheader : 512nB // webp/png/empty 사진헤더.
    common : 4B; KSC5 // KSC5 공통 표지.
    subtype : 4B; KELF // KELF5 표지.
    reserved : 8B; 처음 4바이트는 c0 섹션의 crc32 값입니다.

    data chunk {
        c0 : program info {
            info; str
            abif; int
            sign; nB
            public; str
        }

        c1 : rodata {
            // 스택에 들어갈 상수 값들이 위치.
        }

        c2 : data {
            // 전역변수의 초기화 값들이 위치.
        }

        c3 : text {
            // KVM 바이트코드가 위치.
        }
    }
}

!! 인터럽트 코드는 현재 실행 상태를 나타냅니다. !!
0 : 정상 (계속진행)
1 : 정지 (프로그램 정상종료)
2 : 일반오류 (타입 오류 등, 계속진행 가능)
-1 : 심각오류 (비정상종료)
16 ~ 32 : 테스트 기능 (외부함수)
32 ~ : 실사용 추가기능 (외부함수)
