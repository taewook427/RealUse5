stdlib5.kzip : 파일들/폴더를 단일 파일로 패킹/언패킹.

<python>
class toolbox
    func zipfile(string mode)
    # 파일 -> kzip, 모드는 "png"/"webp"/"nah"
    func zipfolder(string mode)
    # 폴더 -> kzip, 모드는 "png"/"webp"/"nah"
    func unzip(string path)
    # kzip -> 파일/폴더, path 폴더를 초기화 한 후 그 안에 풀림.
    func abs()
    # 내부 파일/폴더 매개변수 표준절대경로화

    .noerr bool # kzip 서브타입과 mainheader crc32 값을 체크하지 않을지 여부. (체크시 False)
    .export string # 결과를 내보낼 폴더/파일 경로. 절대/상대*표준/윈도우 형식 가능.
    .folder string # 패킹할 폴더 경로. 절대*표준/윈도우 형식 가능.
    .file string[] # 패킹할 파일들 경로. 절대*표준/윈도우 형식 가능.

<go>
func Init() -> toolbox
#
struct toolbox
    func Zipfile(string mode)
    # 파일 -> kzip, 모드는 "png"/"webp"/"nah"
    func Zipfolder(string mode)
    # 폴더 -> kzip, 모드는 "png"/"webp"/"nah"
    func Unzip(string path)
    # kzip -> 파일/폴더, path 폴더를 초기화 한 후 그 안에 풀림.
    func Abs()
    # 내부 파일/폴더 매개변수 표준절대경로화

    .Noerr bool # kzip 서브타입과 mainheader crc32 값을 체크하지 않을지 여부. (체크시 False)
    .Export string # 결과를 내보낼 폴더/파일 경로. 절대/상대*표준/윈도우 형식 가능.
    .Folder string # 패킹할 폴더 경로. 절대*표준/윈도우 형식 가능.
    .File string[] # 패킹할 파일들 경로. 절대*표준/윈도우 형식 가능.

!!! 중요 경고 !!!
- 절대 경로는 최초 시작점에서부터의 모든 경로를 나타낸 경로, 상대 경로는 현재 작업 폴더에서의 상대적 경로.
- 표준 경로는 /를 구분자로 사용하는 리눅스 형식 경로, 윈도우 형식은 \를 구분자로 사용하는 윈도우 경로.
- 표준절대경로를 사용하고, 폴더는 뒤에 /를 붙여 쓰는 것을 추천.
- .folder와 .file은 절대경로여야 하며, ../등 상대 경로로 입력했다면 abs()를 사용.
- unzip()의 path와 .export는 절대/상대*표준/윈도우 형식 모두 가능.
- .export 파일/폴더는 결과 출력 시 초기화되며, 폴더 언패킹은 출력 폴더 내부에 원본 파일 모음이 있음.

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
KZIP은 KSC5의 서브타입으로 구성됩니다.
서브타입은 KZIP, 예약 영역은 mainheader의 CRC32 체크섬입니다.
메인헤더는 줄바꿈으로 구분된 폴더/파일 수, 폴더 이름들이 들어갑니다.
모든 문자열 인코딩은 UTF-8로 합니다.

KZIP {
    KSC5 header : nB // 위와 동일한 헤더.
    (
        chunk size (name) : 8B // 이름 청크의 크기.
        chunk (name) : nB // 이름 청크.
        chunk size (file) : 8B // 파일 청크의 크기.
        chunk (file) : nB // 파일 청크. 파일 청크는 0개 혹은 여러 개가 올 수 있다.
    )...
}
청크 개수는 파일 개수 * 2이며, 이름-파일 순서로 반복됩니다.
