stdlib5.kobj [win & linux] : python - golang C FFI 보조기능들.

<python>
func repath() -> str[]
# 현재 실행파일의 작업 디렉토리를 자동으로 맞추고 명령줄 인자를 반환.
func encode(int num, int length) -> bytes
# 정수를 바이트로 리틀 엔디안 인코딩.
func decode(bytes data) -> int
# 리틀 엔디안 인코딩된 바이트를 정수로 디코딩.
func pack(bytes[] series) -> bytes
# 여러 바이트열을 하나의 바이트열로 패키징.
func unpack(bytes chunk) -> bytes[]
# 패키징된 바이트열 덩어리 풀어 바이트열들 반환.
func send(bytes data) -> (CCHAR*, int)
# 바이트열을 ctypes 포인터로 바꾸고 길이 반환.
func recv(CCHAR* reader, int length) -> bytes
# 함수 출력 C char 포인터와 길이를 입력받아 바이트열 반환.
func recvauto(CCHAR* reader) -> bytes
# 함수 출력 C char 포인터를 받아 자동으로 길이 추출해 바이트열 반환.
func call(str args, str rets) -> (tuple, CTYPES)
# py ctypes 함수 입출력 설정값인 튜플, ctypes 타입 값을 반환.

<go>
func Repath() -> str[]
# 현재 실행파일의 작업 디렉토리를 자동으로 맞추고 명령줄 인자를 반환.
func Encode(int num, int length) -> byte[]
# 정수를 바이트로 리틀 엔디안 인코딩.
func Decode(byte[] data) -> int
# 리틀 엔디안 인코딩된 바이트를 정수로 디코딩.
func Pack(byte[][] series) -> byte[]
# 여러 바이트열을 하나의 바이트열로 패키징.
func Unpack(byte[] chunk) -> byte[][]
# 패키징된 바이트열 덩어리 풀어 바이트열들 반환.

<go comment>
func Recv(CCHAR* arr, CINT length) -> byte[]
# C CHAR 포인터와 길이를 받아 바이트 슬라이스 반환.
func Send(byte[] arr) -> CCHAR*
# 바이트 슬라이스를 데이터 그대로 C CHAR 포인터로 변환.
func Sendauto(byte[] arr) -> CCHAR*
# 자동으로 길이헤더를 붙어 C CHAR 포인터로 변환.
func Free(CCHAR* arr)
# malloc으로 할당된 C FFI 메모리를 해제. go -> py로 바이트열을 보낸 후 항상 호출해야 함.

기본적으로 사용하는 모듈이 적기에 FFI가 아닌 다른 목적으로 사용할 수도 있습니다.
go 소스코드에 주석처리된 부분이 GCC-C로 Cgo와 관련된 부분입니다.
주석을 적절히 복사하여 사용하세요.

!! 구조 제한 !!
바이트 배열 패키징 시 255개 이하로만 전송됩니다.
1GiB 이상의 큰 데이터는 전송에 실패할 수도 있습니다.

python에서 go로 제작된 공유 라이브러리를 호출하는것이 전제입니다.
따라서 바이트 배열이 py -> go로 전송될 때 포인터와 배열크기를 모두 전달하고,
포인터에는 데이터 그 자체가 담깁니다. (constB)
반면 go -> py로 갈 때는 데이터 크기가 알려진 경우, 포인터에 데이터 그 자체를 담습니다. (constB)
데이터 크기가 동적인 경우, auto 옵션이 붙은 함수를 사용해야 합니다.
이 함수는 데이터 앞에 크기 8 바이트를 인코딩하여 다양한 크기의 데이터를 지원합니다. (8B + nB)
바이트열 패키징은 1 바이트 개수 + (8 바이트 길이 + n 바이트 데이터) * 반복 구조입니다.
1B len + (8B size + nB data) * n

ctypes.CDLL("./*") 또는 ctypes.cdll.LoadLibrary("./*") 로 오브젝트를 호출합니다.
obj.함수명.argtype / obj.함수명.restype 로 입출력 타입을 정합니다.
python에서 float 전송 시에는 ctypes.c_float 함수로 변환해야 합니다.
python call 함수에서 i : 정수, f : 실수, b : 바이트 배열의 타입을 지원합니다.
