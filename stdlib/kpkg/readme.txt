stdlib5.kpkg [win & linux] : 패키지 설치 & 배포 관리.

<python>
class toolbox
    func pack(int[] osnums, str[] dirpaths, str respath)
    # 지원 운영체제 번호와 패키지 폴더 경로, 출력파일 경로를 입력받아 패키징.
    func unpack(str path) -> str
    # 현재 운영체제에 맞는 패키지를 자동으로 폴더로 언패킹. 언패킹된 폴더 경로 반환.

    .public str # 공개키 텍스트. 미사용시 빈 문자열.
    .private str # 비밀키 텍스트. 미사용시 빈 문자열.
    .name str # 패키지 이름. 공백없이 영문 10글자 내외로 설정 권장.
    .version float # 버전 정보. 숫자가 클수록 최신버전이다.
    .text str # 패키지 설명. 문장 1~2줄 길이로 설정 권장.
    .rel_date str # 배포날짜. 연월일 8글자 문자열.
    .dwn_date str # 설치날짜. 연월일 8글자 문자열.
    .osnum int # 현재 운영체제 정보.

<go>
func Initkpkg(int osnum) -> toolbox
# 운영체제 번호를 설정한 toolbox 구조체 반환.
struct toolbox
    func Pack(int[] osnums, str[] dirpaths, str respath) -> error
    # 지원 운영체제 번호와 패키지 폴더 경로, 출력파일 경로를 입력받아 패키징.
    func Unpack(str path) -> (str, error)
    # 현재 운영체제에 맞는 패키지를 자동으로 폴더로 언패킹. 언패킹된 폴더 경로 반환.

    .Public str # 공개키 텍스트. 미사용시 빈 문자열.
    .Private str # 비밀키 텍스트. 미사용시 빈 문자열.
    .Name str # 패키지 이름. 공백없이 영문 10글자 내외로 설정 권장.
    .Version float # 버전 정보. 숫자가 클수록 최신버전이다.
    .Text str # 패키지 설명. 문장 1~2줄 길이로 설정 권장.
    .Rel_date str # 배포날짜. 연월일 8글자 문자열.
    .Dwn_date str # 설치날짜. 연월일 8글자 문자열.
    .Osnum int # 현재 운영체제 정보.

소프트웨어/확장데이터를 버전, 운영체제에 맞추어 관리하고 설치합니다.
운영체제 번호 (0:모든 OS, 1:Windows, 2:LinuxMint)와 각 패키지 폴더에 따라
하나의 패키지 파일로도 여러 운영체제에서 프로그램을 설치할 수 있습니다.

!!! pack 전 public, private, name, version, text 필드를 설정해야 합니다.
unpack 전 osnum 필드를 설정해야 합니다.
python 버전은 ./temp674, go 버전은 ./temp675 폴더를 임시폴더로 사용하니 작업 후 폴더를 삭제하십시오. !!!

!! 전자서명 사용 시 패키지 데이터를 메모리에 올려야 합니다. 
너무 큰 파일의 경우 전자서명을 끄거나 다른 방식으로 배포하세요. !!

Starter5 프로그램과 연동될 패키지는 두 종류가 있습니다.
Common(공통 기능)/Extension(실행용 프로그램).
모든 패키지는 폴더 형태며, 내부에 _ST5_VERSION.txt 파일이 있습니다.
패키지 설치 시 자동으로 설치 관련 정보가 여기에 기록됩니다.
그 외에도 Extension 패키지에는 _ST5_EXE~가 포함되며,
_ST5_ICON~과  _ST5_DATA/도 패키지에 들어갈 수 있습니다.
