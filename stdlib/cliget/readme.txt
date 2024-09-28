stdlib5.cliget [win & linux] : CLI 환경에서의 다양한 선택창 구현.

<go>
struct PathSel
    func Init(str[] names, str[] paths)
    # PathSel 구조체의 바로가기 경로를 설정함.
    func GetFile(str start) -> string
    # 시작할 폴더를 경로로 입력하고 파일 하나를 선택함.
    func GetFolder(str start) -> string
    # 시작할 폴더를 경로로 입력하고 폴더 하나를 선택함.

struct KeySel
    func Init(byte[] basic, PathSel explorer)
    # KeySel 구조체의 기본 키파일과 파일 탐색기를 설정함.
    func GetKey() -> byte[]
    # 키 파일을 선택해 Data 반환.

    .Data byte[] # 선택된 키 파일 데이터.
    .Path str # 선택된 키 파일 정보.

struct OptSel
    func Init(str[] names, str[] types, str[] limits, PathSel explorer, byte[] basic)
    # 기본 키파일과 파일 탐색기를 설정하고 옵션 이름, 타입, 제한을 설정함. (기본값 설정은 init 직후에 할 것)
    func GetOpt()
    # 옵션 선택 완료 후 종료됨. 옵션 데이터는 .StrRes와 .ByteRes로 확인 가능.

    .Name str[] # 옵션 이름들.
    .StrRes str[] # 문자열 타입 옵션 선택결과.
    .ByteRes byte[][] # 이진 타입 옵션 선택결과.

! 명령어 설명 !
PathSel (tp N) (N) (sel N) (only S) (S)
    바로가기(tp N) : tp와 숫자를 써서 바로가기 폴더로 현재 작업 폴더를 이동한다.
    이동(N) : 상위/하위 폴더로 이동한다. 파일 선택인 경우 파일에 해당하는 숫자라면 파일이 선택된다.
    선택(sel N) : 파일 또는 폴더를 선택한다.
    확장자 설정(only S) : 특정 확장자인 파일만 보여준다. 기본은 전체 모드인 "*"이다. "only *"를 입력해 다시 전체선택 모드로 돌아간다.
    직접입력(S) : 목표 경로를 문자열 그 자체로 직접 입력한다. 앞뒤로 큰따옴표가 있는 "\"경로\"" 형식도 사용 가능하다.
KeySel (submit) (basic) (direct) (comm S)
    제출(submit) : 현재 선택된 키파일 데이터로 제출한다.
    기본으로 설정(basic) : 초기화 때 입력받은 기본키파일로 설정한다.
    직접 경로 입력(direct) : PathSel로 키파일 경로를 입력받는다. 해당 파일의 전체 데이터가 선택된다.
    통신으로 받아오기(comm S) : kcom으로 키파일 정보를 받아온다. 복호화된 데이터가 선택된다.
OptSel (pos content) (submit)
    설정 입력(pos content) : 옵션 번호와 내용을 공백 하나를 두고 입력해 설정한다.
    참/거짓 입력은 t T true True TRUE로, 이진 데이터 입력은 8a9B 같이 16진 표기로 한다.
    제출(submit) : 현재 옵션 설정 상태를 제출한다.

! KeySel kcom 전송 설명 !
direct 모드는 직접 파일을 선택하고 해당 파일을 읽어 반환한다.
이는 키파일이 비암호화 상태로 노출된다는 문제가 있다.
따라서 키파일 데이터를 암호화한 후 kcom으로 데이터를 전송하여
KeySel이 복호화하는 comm 모드가 있다.
comm 모드는 (48B key + nB filepath)가 전송되며,
filepath는 kaes.funcmode로 암호화된 파일 경로이다.

! OptSel 결과저장 설명 !
.Name .StrRes .ByteRes 슬라이스는 길이가 같아 index로 대응되는데,
StrRes와 ByteRes 중 하나에만 옵션 선택 결과가 저장되어 있다.
bool : ByteRes[i]에 True면 0x00 False면 0x01.
int : StrRes[i]에 문자열 형태로 저장.
float : StrRes[i]에 문자열 형태로 저장.
string : StrRes[i]에 값 저장.
bytes : ByteRes[i]에 값 저장.
folder : StrRes[i]에 경로 저장.
file : StrRes[i]에 경로 저장.
keyfile : ByteRes[i]에 데이터 저장.

! OptSel 제한 설명 !
OptSel은 8가지 타입을 인식할 수 있는데, 각 타입마다 설정할 수 있는 값의 제한이 다르다.
bool : "*" 무제한, "T" True만, "F" False만.
int : "*" 무제한, "0" 0만, "0+" 0이상만, "0-" 0이하만, "+" 양수만, "-" 음수만.
float : "*" 무제한, "0" 0과 가까운 값만, "0+" 0이상만, "0-" 0이하만, "+" 양수만, "-" 음수만.
string : "*" 무제한, "N" N바이트 짜리만, "N+" N바이트 이상만, "N-" N바이트 이하만.
bytes : "*" 무제한, "N" N바이트 짜리만, "N+" N바이트 이상만, "N-" N바이트 이하만.
folder : "*" 무제한, "R" 최상위 폴더만, "NR" 최상위 폴더가 아닌것만.
file : "*" 무제한, "ext" 확장자가 ext인것만.
keyfile : "*" 무제한, "N" N바이트 짜리만, "N+" N바이트 이상만, "N-" N바이트 이하만.

CLI 환경에서 옵션 조작을 GUI와 같이 쉽고 직관적으로 하기 위해 제작되었습니다.
표시에 숫자를 사용하며, 입력에 간단한 명령과 숫자만 사용합니다.
모든 구조체는  Init 함수 실행 후에 사용해야 합니다.
