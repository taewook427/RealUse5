stdlib5.ksign [win & linux] : 파일/폴더 통합 해싱, RSA 전자서명 도구.

<python-st>
func khash(str path) -> bytes
# 파일/폴더 통합 64바이트 해시값을 계산.
func kinfo(str path) -> (int size, int file, int folder)
# 파일/폴더에 대해 해당 항목을 포함한 모든 하위 항목의 크기, 파일 개수, 폴더 개수를 계산.
func genkey(int n) -> (str public, str private)
# N bit RSA 키를 생성. N은 1024 이상의 2의 거듭제곱이여야 함.
func sign(str private, bytes plain) -> bytes
# 비밀키를 이용해 특정 바이너리를 서명, 암호화 바이너리를 반환.
func verify(str public, bytes enc, bytes plain) -> bool
# 공개키를 이용해 평문 바이너리, 암호화 바이너리를 입력받아 서명을 검증.

<python-hy>
class toolbox
    func khash(str path) -> bytes
    # 파일/폴더 통합 64바이트 해시값을 계산.
    func kinfo(str path) -> (int size, int file, int folder)
    # 파일/폴더에 대해 해당 항목을 포함한 모든 하위 항목의 크기, 파일 개수, 폴더 개수를 계산.
    func genkey(int n) -> (str public, str private)
    # N bit RSA 키를 생성. N은 1024 이상의 2의 거듭제곱이여야 함.
    func sign(str private, bytes plain) -> bytes
    # 비밀키를 이용해 특정 바이너리를 서명, 암호화 바이너리를 반환.
    func verify(str public, bytes enc, bytes plain) -> bool
    # 공개키를 이용해 평문 바이너리, 암호화 바이너리를 입력받아 서명을 검증.

<go>
func Khash(str path) -> byte[]
# 파일/폴더 통합 64바이트 해시값을 계산.
func Kinfo(str path) -> (int size, int file, int folder)
# 파일/폴더에 대해 해당 항목을 포함한 모든 하위 항목의 크기, 파일 개수, 폴더 개수를 계산.
func Genkey(int n) -> (str public, str private, error)
# N bit RSA 키를 생성. N은 1024 이상의 2의 거듭제곱이여야 함.
func Sign(str private, byte[] plain) -> (byte[], error)
# 비밀키를 이용해 특정 바이너리를 서명, 암호화 바이너리를 반환.
func verify(str public, byte[] enc, byte[] plain) -> (bool, error)
# 공개키를 이용해 평문 바이너리, 암호화 바이너리를 입력받아 서명을 검증.

!!! python VM 경고 !!!
threading 모듈 사용 시 메인 스크립트에
if __name__ == '__main__':
문구를 추가하세요.

!!! py-st 호환 불가 경고 !!!
전자서명의 경우, py-st로 작성된 서명은 go로도 검증할 수 있으나,
go로 작성된 서명을 py-st로 검증할 수 없습니다.
내부적으로 go 런타임을 가진 python-hy 버전을 사용하세요.

khash는 다음 규칙을 따라 64바이트 결과를 내보냅니다.
1. 빈 파일의 해시값은 각 위치가 모두 0인 길이 64의 바이트.
2. 빈 폴더의 해시값은 빈 파일과 동일하게 계산.
3. 비지 않은 파일의 해시값은 파일 바이너리를 SHA3-512 해싱한 결과물.
4. 비지 않은 폴더의 해시값은 다음과 같이 계산함.
    4-1. 현재 폴더의 위치를 정규화 (리눅스 형식, 뒤가 /로 끝나게)
    4-2. 현재 폴더의 하위 항목의 이름에 파일은 그대로, 폴더는 뒤에 /을 붙임.
    4-3. 파일/폴더 상관없이 UTF-8 인코딩한 이름의 바이트순에 따라 오름차순 정렬함.
    4-4. 각 항목의 해시값 64바이트를 정렬된 순으로 이어붙여 하나의 바이너리를 만듦.
    4-5. 폴더의 해시값은 이 바이너리를 SHA3-512 해싱한 결과물.

전자서명의 저장 양식은 PKIX (public), PKCS1 (private), PEM, PSS 입니다.
