# test561 : kweb (py)

import requests
import bs4

# http~.html에서 domain 별 txt 반환
def gettxt(url, domain):
    html = requests.get(url).text
    dom = bs4.BeautifulSoup(html, "html.parser")
    data = dom.find("p", id = domain)
    if data == None:
        raise Exception(f"No Domain : {domain}")
    else:
        return data.text

# http~/ + name + .num 를 path로 바이너리 생성
def download(url, name, num, path):
    if url[-1] != "/":
        url = url + "/"
    with open(path, "wb") as f:
        for i in range(0, num):
            html = requests.get(url + name + f".{i}")
            if html.content == None:
                raise Exception(f"No Data : {name}.{i}")
            else:
                f.write(html.content)

# edgedriver 홈페이지에서 정보 받아오기, 홈페이지 + 행동정보 + zip경로 -> 버전 정수
def driver(url, datas, path):
    # html 요소, 공통 레이어 class 명, html 요소, 버튼 class 명, html 요소, 버전 class 명
    # ["div", "common-card-list__card", "a", "common-button.common-button--tag", "div", "block-web-driver__versions"]
    html = requests.get(url).text
    dom = bs4.BeautifulSoup(html, "html.parser")
    cand = dom.findAll( datas[0], { "class" : datas[1] } )
    box = None
    for i in cand:
        if "stable" in i.text.lower().replace(" ", ""):
            box = i
    if box == None:
        raise Exception("No Stable WD")
    but = None
    cand = box.findAll( datas[2], { "class" : datas[3] } )
    for i in cand:
        if "x64" in i.text.lower().replace(" ", ""):
            but = i
    if but == None:
        raise Exception("No x64 WD")
    link = but.attrs["href"] # x64 driver link
    cand = box.find( datas[4], { "class" : datas[5] } )
    if cand == None:
        ver = 0
    else:
        ver = ""
        num = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9"]
        for i in cand.text:
            if i in num:
                ver = ver + i
            elif i == "." and num != "":
                break
        ver = int(ver) # 버전 정수
    with open(path, "wb") as f:
        html = requests.get(link)
        if html.content == None:
            raise Exception(f"No Data : {link}")
        else:
            f.write(html.content)
    return ver

# CM이 관리하는 사용자 영역은 extension과 common
# binary와 value는 Starter5와 CM이 자동 관리.
# extension : CM 제외 extension
# common : webdriver 제외 공통 프로그램
# binary : webdriver, CM
# value : MS web data, values -> CM이 common에 자동 업로드
"""
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
    0.num = 1; 0.name = "msedgedriver.exe"; 0.version = 120.0; 0.data = 20240112;<br>
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
"""
