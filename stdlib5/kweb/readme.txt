stdlib5.kweb : gen5 datasvr 받아오기, MS edgedriver 다운로드

<python>
func gettxt(str url, str domain) -> str kdbtxt
# http~.html 주소와 그 안의 id domain을 받아 kformat str 반환.
func download(str url, str name, int num, str path)
# http~/ 주소와 파일 이름(.N 전), 조각 개수 받아 경로에 원본 파일 생성.
func driver(str url, str[] datas, str path) -> int ver
# MS edgedriver 홈페이지 주소와 위치 파라미터를 받아 경로에 driver zip 파일 생성.

<go>
func Gettxt(str url, str domain) -> str kdbtxt
# http~.html 주소와 그 안의 id domain을 받아 kformat str 반환.
func Download(str url, str name, int num, str path)
# http~/ 주소와 파일 이름(.N 전), 조각 개수 받아 경로에 원본 파일 생성.
func Driver(str url, str[] datas, str path) -> int ver
# MS edgedriver 홈페이지 주소와 위치 파라미터를 받아 경로에 driver zip 파일 생성.

kweb 제작 기준 파라미터
https://taewook427.github.io/lite-web/gen5_datasvr/test563.html
extension common binary value
https://taewook427.github.io/lite-web/gen5_datasvr/
~
https://developer.microsoft.com/ko-kr/microsoft-edge/tools/webdriver/
["div", "common-card-list__card", "a", "common-button.common-button--tag", "div", "block-web-driver__versions"]
~.zip

CM이 관리하는 사용자 영역은 extension과 common.
binary와 value는 Starter5와 CM이 자동 관리.
extension : CM 제외 extension.
common : webdriver 제외 공통 프로그램.
binary : webdriver, CM.
value : MS web data, values -> CM이 common에서 자동 업데이트.

<extension>
# Starter5 확장, 사용자가 CM으로 관리 가능.
.num .name .txt

<common>
# Starter5 공통, 사용자가 CM으로 관리 가능.
.num .name .txt

<binary>
# 내부 데이터 공유용, MSWD 백업본 위치.
.num .name .version

<value>
# 내부 값 공유용.
(필드 없이 순수 kformat 텍스트.)

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>kweb datasvr 0</title>
</head>
<body>
    <p id="extension">
    num = 1;<br>
    0.num = 3; 0.name = "gen5win_reltest.zip"; 0.txt = "윈도우용 테스트 프로그램";<br>
    </p>

    <p id="common">
    num = 1;<br>
    0.num = 3; 0.name = "gen5win_reltest.zip"; 0.txt = "윈도우용 테스트 프로그램";<br>
    </p>

    <p id="binary">
    num = 1;<br>
    0.num = 1; 0.name = "msedgedriver.exe"; 0.version = 120;<br>
    </p>
    
    <p id="value">
    name = "Ko";
    <br>
    age = 21;
    </p>

    <p id="explain">
    extension과 common은 content manager용.<br>
    CM에서 binary와 value 데이터 업데이트 가능.<br>
    특별히 webdriver는 MS 홈페이지와 binary로 자동 업데이트.<br>
    binary에 CM은 Starter5용.
    </p>
    
</body>
</html>
