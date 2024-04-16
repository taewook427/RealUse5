stdlib5.ksc [win & linux] : 다양한 상황에 사용할 수 있는 파일 청크 표준형.

<python>
class toolbox
    func readf()
    # 파일을 읽어 내부 저장소 정보 업데이트.
    func writef()
    # 내부 저장소 정보 기반 새 KSC5 파일 생성.
    func addf(str path)
    # 주어진 경로의 파일을 KSC5의 청크 하나로 추가 기록.
    func linkf(bytes data)
    # 주어진 바이너리를 KSC5의 청크 하나로 추가 기록.
    func readb(bytes data)
    # 주어진 바이너리로 내부 저장소 정보 업데이트.
    func writeb() -> bytes
    # 내부 저장소 정보 기반 새 KSC5 바이너리 생성.
    func addb(bytes stream, str path) -> bytes
    # 주어진 경로의 파일을 KSC5 스트림의 끝에 청크 하나로 추가 기록, 바이너리 반환.
    func linkb(bytes stream, bytes data) -> bytes
    # 주어진 바이너리를 KSC5 스트림의 끝에 청크 하나로 추가 기록, 바이너리 반환.

    .prehead bytes # 메인헤더 앞의 가짜헤더. (512nB) 기본값은 webp 512바이트.
    .common bytes # KSC5 공통 식별자. (4B; KSC5)
    .subtype bytes # 하위 타입 식별자. (4B)
    .reserved bytes # 프로그램별 다목적 예약부분. (8B)
    .headp int # 메인헤더 시작 오프셋 (512nB, prehead 길이와 동일) 기본값은 512.
    .rsize int # 실질 데이터 크기 (메인헤더 + 데이터청크 길이) 기본값은 0.
    .path str # 파일 입출력용 경로 지정 문자열.
    .predetect bool # 데이터 청크 정보를 미리 읽어올지 여부. 기본값은 False.
    .chunkpos int[] # 데이터 청크 (8B + nB)의 시작 오프셋 정보.
    .chunksize int[] # 데이터 청크 (8B + nB)의 크기 정보 (8 + n).

func encode(int num, int length) -> bytes
# 정수를 바이트로 리틀 엔디안 인코딩.
func decode(bytes data) -> int
# 리틀 엔디안 인코딩된 바이트를 정수로 디코딩.
func crc32hash(bytes data) -> bytes
# 입력 데이터의 CRC32 값을 반환.
func webpbase() -> bytes
# KSC5 기본 prehead 512B (webp 형태) 반환.

<go>
func Initksc() -> toolbox
# toolbox 구조체의 초기값을 설정한 후 반환.
struct toolbox
    func Readf() -> error
    # 파일을 읽어 내부 저장소 정보 업데이트.
    func Writef() -> error
    # 내부 저장소 정보 기반 새 KSC5 파일 생성.
    func Addf(str path) -> error
    # 주어진 경로의 파일을 KSC5의 청크 하나로 추가 기록.
    func Linkf(byte[] data) -> error
    # 주어진 바이너리를 KSC5의 청크 하나로 추가 기록.
    func Readb(byte[] data) -> error
    # 주어진 바이너리로 내부 저장소 정보 업데이트.
    func Writeb() -> (byte[], error)
    # 내부 저장소 정보 기반 새 KSC5 바이너리 생성.
    func Addb(byte[] stream, str path) -> (byte[], error)
    # 주어진 경로의 파일을 KSC5 스트림의 끝에 청크 하나로 추가 기록, 바이너리 반환.
    func linkb(byte[] stream, byte[] data) -> byte[]
    # 주어진 바이너리를 KSC5 스트림의 끝에 청크 하나로 추가 기록, 바이너리 반환.

    .Prehead byte[] # 메인헤더 앞의 가짜헤더. (512nB) 기본값은 webp 512바이트.
    .Common byte[] # KSC5 공통 식별자. (4B; KSC5)
    .Subtype byte[] # 하위 타입 식별자. (4B)
    .Reserved byte[] # 프로그램별 다목적 예약부분. (8B)
    .Headp int # 메인헤더 시작 오프셋 (512nB, Prehead 길이와 동일) 기본값은 512.
    .Rsize int # 실질 데이터 크기 (메인헤더 + 데이터청크 길이) 기본값은 0.
    .Path str # 파일 입출력용 경로 지정 문자열.
    .Predetect bool # 데이터 청크 정보를 미리 읽어올지 여부. 기본값은 false.
    .Chunkpos int[] # 데이터 청크 (8B + nB)의 시작 오프셋 정보.
    .Chunksize int[] # 데이터 청크 (8B + nB)의 크기 정보 (8 + n).

func Encode(int num, int length) -> byte[]
# 정수를 바이트로 리틀 엔디안 인코딩.
func Decode(byte[] data) -> int
# 리틀 엔디안 인코딩된 바이트를 정수로 디코딩.
func Crc32hash(byte[] data) -> byte[]
# 입력 데이터의 CRC32 값을 반환.
func Webpbase() -> byte[]
# KSC5 기본 prehead 512B (webp 형태) 반환.

다양한 상황에서 사용할 수 있는 데이터 모음 형식입니다.
사진 데이터를 포함시켜 아이콘을 위장시킬 수 있는 pre-header와
식별자와 세부 타입을 포함한 main-header 뒤에
데이터를 담은 데이터 청크가 0개 이상 반복됩니다.
만약 청크의 종결을 확인할 수 있다면, 그 이후에 오는 값들은 trash가 됩니다.

KSC5 {
    pre-header + padding : 512nB // 사진 정보를 담을 수 있는 가짜 헤더, 0또는 512B 배수 크기.
    common sign : 4B; KSC5 // KSC5 파일임을 나타냄.
    subtype sign : 4B // KSC5 파일 중 어느 종류인지 나타냄.
    reserved : 8B // 다양한 방법으로 사용할 수 있게 예약된 영역.
    data chunk {
        chunk size : 8B // 리틀 엔디안 인코딩된 청크의 크기.
        chunk data : nB // 청크 데이터.
    } * n // 데이터 청크는 0개 이상 반복해서, 청크 사이 빈틈 없이 나타남.
    trash : nB // 청크가 모두 끝나고 남은 부분. (쓰레기값)
}

predetect 옵션이 참인경우, read 단계에서 모든 청크의 위치와 크기를 미리 리스트에 기록합니다.
이 정보는 청크 사이즈와 데이터 부분을 합한것을 기준으로, chunk size 필드의 시작 부분의 위치를 기록합니다.
청크 사이즈 8 바이트의 모든 값이 255 (8x FF)인 경우, 청크 종결표현으로 인식됩니다.
이 부분은 청크 데이터로 취급되지 않으며, 이후에 오는 모든 데이터는 쓰레기값입니다.

!! 청크 종결표현 삽입 !!
add 함수에 파일 경로가 아닌 빈 문자열을 넣을 경우 종결표현이 기록됩니다.
종결표현의 삽입 여부는 필수가 아닌 선택입니다.
