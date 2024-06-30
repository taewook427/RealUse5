stdlib5.kgui [win & linux] : 크로스 플렛폼 GUI 템플릿.

<python linesel>
class toolbox(str title, bool iswin)
    func entry()
    # 화면 초기화. curpos 설정 후 호출해야 함.
    func render()
    # 화면 그리기. curpos에 맞게 항목이 보여짐.
    func guiloop()
    # GUI 메인루프. 실행 시 GUI 창이 꺼질 때까지 루프를 돈다.
    func custom0(int x, int y)
    # options[x][y] 항목이 더블클릭 되었을 때 실행되는 함수 틀.

    .infos str[] # 각 대분류 별 이름 목록.
    .options str[][] # 대분류-소분류 세분화된 선택지 항목들.
    .mwin TkObj # 메인 윈도우 창 TkInter 객체.
    .curpos int # 현재 보여지는 대분류 번호.

<python blocksel>
class toolbox(str title, bool iswin)
    func entry()
    # 화면 초기화. upos, curpos 설정 후 호출해야 함.
    func render(bool uonly)
    # 화면 그리기. upos, curpos에 맞게 항목이 보여짐. uonly가 True면 위쪽 메세지 표시 부분만 새로 그려짐.
    func guiloop()
    # GUI 메인루프. 실행 시 GUI 창이 꺼질 때까지 루프를 돈다.
    func custom0(int x)
    # txts[x] 항목이 클릭 되었을 때 실행되는 함수 틀.

    .pics str[] # 선택항목 사진파일 경로. 150x150 크기여야 함.
    .txts str[] # 선택항목 이름. 영문 12글자 이내여야 함.
    .umsg str[] # 상단 메세지 표시 라벨에 보여질 메세지들.
    .mwin TkObj # 메인 윈도우 창 TkInter 객체.
    .upos int # 현재 보여지는 상단 메세지 번호.
    .curpos int # 현재 보여지는 선택 갤러리.

<python explorer>
class toolbox(str title, bool iswin)
    func entry()
    # 화면 초기화. icons, tabpos, viewpos 설정 후 호출해야 함.
    func render(bool upper, bool view0, bool view1, bool sch, bool txtd, bool picd, bool bind)
    # 화면 그리기. 각 불리언 인자가 참일때만 해당 부분/탭이 업데이트됨.
    # upper : 상단 주소표시줄, view 소그룹 표시. view0 : view 탭 아이콘사진, 이름, 크기. view1 : view 탭 잠김여부, 선택여부. sch/txtd/picd/bind : 각 탭 전체.
    func guiloop()
    # GUI 메인루프. 실행 시 GUI 창이 꺼질 때까지 루프를 돈다. 0.1초에 한번씩 custom9 함수를 호출하고 참이면 search 탭 새로고침 활성화 표시.
    func menubuilder()
    # 메뉴 바 생성 함수 틀. entry 함수 호출 시 자동호출됨. 이 함수를 상속구현하여 mbar 필드에 메뉴 바와 동작 함수를 입힌다.
    func custom0(x)
    # menu option[x] 항목 선택 시 실행되는 함수 틀. !! 메뉴 바 함수로 직접 바인딩해야 한다. !!
    func custom1(x)
    # 상단버튼 클릭 시 실행되는 함수 틀. 상단 확인 버튼(x=0), 상단 뒤로가기 버튼(x=1).
    func custom2(x)
    # viewA(사진 아이콘) 버튼 클릭 시 실행되는 함수 틀. names[x]로 접근가능할때만 호출됨.
    func custom3(x)
    # viewB(잠금여부) 버튼 클릭 시 실행되는 함수 틀. names[x]로 접근가능할때만 호출됨.
    func custom4(x)
    # viewC(선택여부) 버튼 클릭 시 실행되는 함수 틀. names[x]로 접근가능할때만 호출됨.
    func custom5(x, y)
    # search 탭 버튼 클릭 시 실행되는 함수 틀. 검색(x=0), 새로고침(x=1). y는 검색상자 입력 문자열.
    func custom6(x)
    # txt edit 탭 save 버튼 클릭 시 실행되는 함수 틀. x는 텍스트박스 입력 문자열.
    func custom7(x)
    # pic view 탭 버튼 클릭 시 실행되는 함수 틀. 이전항목이동(x=0), 다음항목이동(x=1).
    func custom8(x)
    # bin edit 탭 save 버튼 클릭 시 실행되는 함수 틀. x는 입력된 바이트열. 
    func custom9() -> bool
    # 0.1초에 한번씩 호출되며, 명시적 새로고침 활성화 표시를 원할 시에만 참을 반환하도록 해야하는 함수 틀.

    .paths str[5] # 각 탭마다 상단 주소표시줄에 보여질 내용들.
    .names str[] (size N) # view 탭에서 보여질 항목들의 이름들. 폴더는 /로 끝나야 함.
    .sizes int[] (size N) # view 탭에서 보여질 항목들의 크기들. 음수는 정보없음으로 표시됨.
    .locked bool[] (size N) # view 탭에서 보여질 항목들의 잠김여부. True가 잠김(초록색)으로 표시됨.
    .selected bool[] (size N) # view 탭에서 보여질 항목들의 선택여부.  True가 선택됨(초록색)으로 표시됨.
    .search str[] # search 탭에서 보여질 검색결과들.
    .log0 str # search 탭에서 보여질 현재상태.
              # search 탭 우측 새로고침 버튼을 활성화(초록색)시켜 로그데이터의 변화를 알릴 수 있음.
    .log1 str[] # search 탭에서 보여질 활동로그.
    .txtdata str # txt edit 탭에서 보거나 수정될 텍스트 데이터.
    .picdata bytes # pic view 탭에서 보여질 사진의 이진 데이터.
    .bindata bytes # bin edit 탭에서 보거나 수정될 2GiB 미만의 이진 데이터.
    .icons (bytes IconPic)[str type] # 사진의 크기는 100x100으로 자동조절됨.
           # view 탭 항목 확장자에 해당되는 아이콘 사진 바이너리.
           # "none", "file", "dir" 키는 반드시 존재해야 하며, 각 확장자에 해당하는 사진의 키는 ".ext" 소문자 형태.
    .taboos int (0~4) # 현재 보고 있는 탭의 번호.
    .viewpos int (1~) # 현재 보고 있는 view 탭 소그룹의 번호. (1부터 시작하며 한번에 12개씩 보여짐.)
    .mwin TkObj # 메인 윈도우 창 TkInter 객체.
    .mbar TkObj # 메뉴 바 TkInter 객체.

kgui는 3개의 python 코드 파일을 제공한다. 사용하고자 하는 GUI 용도에 맞게 파일 중 하나를 선택해 import한다.
코드 사용 시 라이브러리의 클래스를 상속받고 customN 함수를 구현하여 상호작용한다.
상위 클래스 생성 시 창 제목과 Windows OS를 쓰는지 여부를 입력받아야 한다.
super().__init__(str, bool)로 초기 설정을 한 후 클래스 내부 값을 설정한다.
그 후 entry()로 창을 초기화 한 후, 클래스 내부의 현재 보기 정보를 설정한다.
마지막으로 guiloop() 함수로 루프를 돌며 사용자 상호작용을 처리하면 된다.

linesel은 대량의 선택지를 효과적으로 고르는 상황에 사용하도록 디자인되었다.
텍스트와 리스트박스 조합으로 더블클릭을 통해 상호작용한다.
아래쪽 버튼으로 다른 대분류로 이동할 수 있다.
사용 순서 : .infos/.options -> entry -> .curpos ->  guiloop

blocksel은 10개 내외의 선택지를 그림표지와 함께 고르는 상황에 사용되도록 디자인되었다.
상단 버튼으로 보여지는 메세지 내용과 upos 값을 바꿀 수 있다.
한번에 6개 선택지씩 보여지며, pics와 txts 항목의 수는 1 이상인 6의 배수여야 한다.
양 옆 버튼으로 다른 선택지가 보여지도록 할 수 있다.
사용 순서 : .pics/.txts/.umsg -> entry -> .upos/.curpos -> guiloop

explorer는 간단한 파일시스템 조작을 위한 다양한 기능을 갖추도록 디자인되었다.
각 탭은 파일탐색기, 로그뷰어, 텍스트에디터, 사진뷰어, 바이너리에디터 기능을 담당한다.
인터페이스 설계사상은 다음과 같다. 모든 순간마다 어떤 "상태"가 있고, 사용자는 이 상태를 그래픽을 통해 상호작용한다.
사용자의 행위 중 "상태"에 변화를 주는 것은 상속구현 함수를 호출한다. 각 함수들은 인터페이스만을 담당하며 실제 업무는 다중스레드상에서 동작한다.
사용 순서 : .icons/.tabpos/.viewpos -> entry -> .paths/.names/.sizes/.locked/.selected/.search/.log0/.log1/.txtdata/.picdata/.bindata -> guiloop
(상황에 따라 필드값을 바꾸며 render 함수를 사용한다.)
