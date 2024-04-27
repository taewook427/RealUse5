stdlib5.kcom [win & linux] : 인터넷 HTTP / 컴퓨터 socket 통신 기능.

<py>
func gettxt(str url, str domain) -> str
# http~.html 주소와 그 안의 id domain을 받아 kformat str 반환.
func download(str url, str name, int num, str path, list proc)
# http~/ 주소와 파일 이름(.N 전), 조각 개수 받아 경로에 원본 파일 생성.
func pack(int port, bytes key) -> str
# 포트 번호와 세션 키 4 바이트를 문자열로 패키징.
func unpack(str address) -> (int, bytes)
# 패키징된 문자열을 포트 번호와 세션 키 4 바이트로 언패킹.

class node
    func send(bytes data, bytes key)
    # 이진 데이터와 세션 키 4 바이트를 입력받아 소켓 전송. 전송 후 소켓은 닫힘.
    func recieve(bytes key) -> bytes
    # 세션 키 4 바이트를 입력받아 이진 데이터 수신. 수신 후 소켓은 닫힘.

    bool .ipv6 # IPv6 주소를 사용할지 여부. 기본값은 True(IPv6사용).
    int .port # 포트 번호. 기본값은 13600번 포트.
    int .close # 타임아웃 시간. 기본값은 150초. 음수일 경우 시간초과 오류 없음.

<go>
func Gettxt(str url, str domain) -> (str, error)
# http~.html 주소와 그 안의 id domain을 받아 kformat str 반환.
func Download(str url, str name, int num, str path, float* proc) -> error
# http~/ 주소와 파일 이름(.N 전), 조각 개수 받아 경로에 원본 파일 생성.
func Pack(int port, byte[] key) -> str
# 포트 번호와 세션 키 4 바이트를 문자열로 패키징.
func Unpack(str address) -> (int, byte[], error)
# 패키징된 문자열을 포트 번호와 세션 키 4 바이트로 언패킹.

func Initcom() -> node
# 송수신 노드 구조체의 초기값을 설정.
struct node
    func Send(byte[] data, byte[] key) -> error
    # 이진 데이터와 세션 키 4 바이트를 입력받아 소켓 전송. 전송 후 소켓은 닫힘.
    func Recieve(byte[] key) -> (byte[], error)
    # 세션 키 4 바이트를 입력받아 이진 데이터 수신. 수신 후 소켓은 닫힘.

    bool .Ipv6 # IPv6 주소를 사용할지 여부. 기본값은 True(IPv6사용).
    int .Port # 포트 번호. 기본값은 13600번 포트.
    int .Close # 타임아웃 시간. 기본값은 150초. 음수일 경우 시간초과 오류 없음.

!! 운영체제 호환성 주의 !!
윈도우와 리눅스는 소스코드 파라미터가 약간 달라서, 직접 수정해야 합니다.

!!! 프로세스 진행도 (proc) !!!
입력받는 파라미터의 위치에 진행도를 기록합니다.
py는 [-1.0], go는 &(-1.0) 등의 방식으로 사용하세요.

!!! 중요 경고 !!!
KCOM5 통신은 동기적으로 동작합니다. 데이터를 모두 전송하거나 응답을 모두 받을 때까지 함수가 반환되지 않습니다.
서버 소켓이 열리지 않은 상태에서는 클라이언트 소켓을 열 수 없습니다.
한 쌍의 프로세스끼리만 소량의 바이너리 데이터가 통신 가능합니다. (일대일 통신)
1023 이하 포트는 예약된 특수 포트이니 10000 ~ 40000 대역의 임의의 포트 사용을 권장합니다.
기본적으로 4B CRC / 4B encryption을 제공하나, 더 안정적인 통신이나 보안을 위해 추가 레이어를 거치십시오.

프로세스 간 통신을 위해 TCP 소켓을 사용합니다.
python과 go 사이 상호 운용이 가능합니다.
다음과 같은 방식으로 서버-클라이언트가 통신하며 서버측의 데이터를 클라이언트에게 보냅니다.
1. 서버 소켓이 열립니다.
1b. 만약 시간초과 설정값이 양수이고 지정 시간 동안 클라이언트 연결이 없다면 타임아웃으로 종료됩니다.
2. 클라이언트 소켓이 열립니다.
3. 클라이언트 측에서 kcom5 + 3바이트 난수를 보냅니다.
4. 서버 측에서 8바이트 값을 읽고 kcom5로 시작하는지 확인합니다.
4b. 만약 kcom5로 시작하지 않는 값을 보내왔다면 잘못된 연결로 종료합니다.
5. 서버 측에서 8바이트 값을 다시 반송합니다.
6. 클라이언트 측에서 8바이트 값을 읽고 보낸 값과 동일한지 확인합니다.
6b. 만약 동일한 8바이트 값이 아니라면 잘못된 연결로 종료합니다.
7. 서버 측에서 메세지 길이(8B) + 원본 메세지 CRC(4B) + 치환된 메세지를 보냅니다.
8. 소켓이 닫히고 서버 측에서 통신이 종료됩니다.
8b. 클라이언트 쪽에서 메세지를 복호화하고 CRC 값을 체크합니다.

기존 stdlib5.kcom / stdlib5.kweb의 기능을 통합했습니다.
서버 html로부터 특정 id의 p 텍스트를 받아오거나, *.숫자 형식의 데이터를 다운로드하는 기능은 동일합니다.
EdgeWD 업데이트와 같은 기능은 다른 전용 프로그램을 통해 제공될 예정입니다.
