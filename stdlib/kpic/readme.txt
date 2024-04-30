stdlib5.kpic [win & linux] : 데이터 사진화.

<python>
class toolbox
    func setmold(str path, int row, int col)
    # 인코딩 시 사용할 사진을 설정합니다. path로 빈 문자열을 전달하면 picdt.kpic5 기본 사진이 사용됩니다.
    # row, col은 인코딩 결과물의 가로, 세로 크기이며 둘 중 하나라도 -1일 경우 사진 크기 그대로 설정됩니다.
    # row, col 값은 모두 4의 배수여야 합니다.
    func detect() -> (str name, int num, str style)
    # target 폴더 안의 항목에서 kpic 사진들을 감지해 정보를 추출합니다.
    # target 설정 후 사용하세요.
    func pack(int zmode) -> (str name, int num)
    # 데이터 파일을 사진으로 인코딩합니다. zmode는 파일 1 바이트당 대응 픽셀 바이트 수 입니다. (2 or 4)
    # setmold 함수를 먼저 사용한 상태여야 합니다. target, export, style 설정 후에 사용하세요.
    func unpack(str name, int num)
    # 사진을 데이터 파일로 디코딩합니다.
    # target, export, style 설정 후에 사용하세요.
    func restore(str[] files, int zmode) -> (str name, int maxnum)
    # 사진의 이름을 내부 데이터로부터 추출해 복구합니다.
    # style 설정 후에 사용하세요.

    .target str # 읽기 대상이 되는 폴더/파일 경로.
    .export str # 쓰기 대상이 되는 폴더/파일 경로.
    .style str # 사진 형식. "webp", "png", "bmp" 지원.
    .proc float # 진행도 변수. (-1.0 : 시작 전, 0.0 ~ 1.0 : 작업 중, 2.0 : 완료)

<go-exp>
func Initpic(str path, int row, int col) -> (toolbox, error)
# 사진 경로를 설정한 toolbox 구조체를 반환합니다. path로 빈 문자열을 전달하면 picdt.kpic5 기본 사진이 사용됩니다.
# row, col은 인코딩 결과물의 가로, 세로 크기이며 둘 중 하나라도 -1일 경우 사진 크기 그대로 설정됩니다.
# row, col 값은 모두 4의 배수여야 합니다.
struct toolbox
    func Detect() -> (str name, int num, str style)
    # Target 폴더 안의 항목에서 kpic 사진들을 감지해 정보를 추출합니다.
    # Target 설정 후 사용하세요.
    func Pack(int zmode) -> (str name, int num)
    # 데이터 파일을 사진으로 인코딩합니다. zmode는 파일 1 바이트당 대응 픽셀 바이트 수 입니다. (2 or 4)
    # Target, Export, Style 설정 후에 사용하세요.
    func Unpack(str name, int num)
    # 사진을 데이터 파일로 디코딩합니다.
    # Target, Export, Style 설정 후에 사용하세요.
    func Restore(str[] files, int zmode) -> (str name, int maxnum)
    # 사진의 이름을 내부 데이터로부터 추출해 복구합니다.
    # Style 설정 후에 사용하세요.

    .Target str # 읽기 대상이 되는 폴더/파일 경로.
    .Export str # 쓰기 대상이 되는 폴더/파일 경로.
    .Style str # 사진 형식. "webp", "png", "bmp" 지원.
    .Proc float # 진행도 변수. (-1.0 : 시작 전, 0.0 ~ 1.0 : 작업 중, 2.0 : 완료)

! go-exp !
golang은 현재 안정적인 webp 인코더를 찾을 수 없는 관계로 experimental 버전만 제공됩니다.
webp style의 pack 진행 시 libwebp 실행파일이 필요합니다.
OS에 맞춰 cwebp 실행파일명을 지정하십시오.

!!! python VM 경고 !!!
multiprocessing 모듈 사용 시 메인 스크립트에
if __name__ == '__main__':
문구를 추가하세요.

!!! 중요 경고 !!!
- pack 함수는 출력 경로의 폴더를 자동으로 초기화합니다. 출력 폴더 내부에 데이터가 있게 하지 마세요.
- go-exp는 webp 인코딩 시 프로그램이 먼저 끝나버리면 cwebp 프로세스가 종료되어버립니다.
- 사진 크기에 맞춰 패딩 바이트가 추가되며, 파일 이름은 저장되지 않으니 kzip과 같이 운용하세요.

사진 바이트 속에 이진 데이터를 숨길 수 있습니다.
기본 원리는 사진의 픽셀 바이트에서 16으로 나눈 몫은 그대로,
나머지를 이용해 바이트를 인코딩하는 것입니다.
모드 2에서는 1바이트 인코딩에 사진 2바이트가, 모드 4에서는 사진 4바이트가 쓰입니다.
사용자 데이터는 사진 크기에 맞춰 분할되며, 각 사진 속 데이터의 첫 8 바이트는
일련번호, 현재 번호, 전체 개수입니다.
기본으로 제공되는 little.webp, standard.webp, big.webp는
모드 2의 경우 각각 1.318 MiB, 9.670 MiB, 47.461 MiB까지 저장할 수 있습니다.
