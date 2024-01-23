stdlib5.kpic : 데이터 사진화

<python>
class toolbox
    func setmold(string path)
    # 주형 사진을 설정.
    func convmold(int row, int col)
    # 주형 사진을 사용 할 때의 해상도를 설정. 임시폴더 초기화.
    func clear(bool mknew)
    # 임시 폴더를 삭제, mknew가 True면 새로 만듦.
    func pack(int mode) -> string name, int num
    # 타겟 파일을 스타일대로 변환해 출력 폴더에 생성, mode는 2/4.
    func unpack(string name, int num)
    # 타겟 폴더의 사진들을 출력 경로에 파일로 생성.
    func detect() -> string name, int num, string style
    # 타겟 폴더에 kpic 사진이 존재하는지, 일련번호와 전체 개수 반환.
    func restore(string[] files, int mode) -> int num
    # 사진에 내장된 고유번호로 원래 이름 복구.

    .moldsize [int row, int col] # 주형 사진의 가로 * 세로 크기.
    .moldpath string # 주형 사진의 경로.
    .temppath string # 임시 폴더 경로. "./temp574/" 사용을 권장.
    .target string # 변환 대상인 파일/폴더.
    .export string # 출력 장소인 파일/폴더.
    .style string # 변환 대상/출력 대상의 확장자.

!!! 중요 경고 !!!
- pack은 임시폴더가 초기화되어야 합니다. set/convmold를 사용한 후 쓰세요.
- pack은 target 파일, export 폴더, style를 설정한 후 써야 합니다.
- unpack은 target 폴더, export 파일, style 설정 후 쓰세요.
- detect은 target 폴더 설정이 필요합니다.
- restore은 파일들 이름, 모드, export 폴더 설정이 필요합니다.
- style은 "bmp", "png", "webp", mode는 2, 4가 가능합니다.
- 사진 크기에 맞춰 패딩 바이트가 추가되며, 파일 이름은 저장되지 않으니 kzip과 같이 운용하세요.

!!! python VM 경고 !!!
multiprocessing 모듈 사용 시 메인 스크립트에
import multiprocessing as mp
if __name__ == '__main__':
    mp.freeze_support()
문구를 추가할 것.

사진 바이트 속에 이진 데이터를 숨길 수 있습니다.
기본 원리는 사진의 픽셀 바이트에서 16으로 나눈 몫은 그대로,
나머지를 이용해 바이트를 인코딩하는 것입니다.
모드 2에서는 1바이트 인코딩에 사진 2바이트가, 모드 4에서는 사진 4바이트가 쓰입니다.
사용자 데이터는 사진 크기에 맞춰 분할되며, 각 사진 속 데이터의 첫 8 바이트는
일련번호, 현재 번호, 전체 개수입니다.
기본으로 제공되는 little.webp, standard.webp, big.webp는
모드 2의 경우 각각 1.318 MiB, 9.670 MiB, 47.461 MiB까지 저장할 수 있습니다.
현재 golang의 경우는 webp의 안정적인 변환 모듈을 찾지 못한 관계로 제작하지 않습니다.
