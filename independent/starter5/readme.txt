independent.starter5 [win & linux] : 실행 매니저.

<Package Execute>
패키지 프로그램을 실행.

<Open Folder>
패키지 폴더 열기.

<Open Data Storage>
패키지 데이터 폴더 열기. 데이터 폴더는 패키지 업데이트를 해도 내용이 유지된다.

<Show Info>
패키지 정보 보기.

실행 매니저는 다른 프로그램 패키지를 관리하고 실행할 수 있게 모아 보여주는 도구입니다.
extension엔 프로그램 패키지가 들어가며 common엔 공통 라이브러리가 들어갑니다.
직접 수동으로 패키지를 관리하거나 설정을 조작할 수도 있지만,
따로 제공되는 패키지 관리 프로그램을 사용하는 것을 권장합니다.

kpkg를 사용한 패키지 관리 프로그램은 따로 제공됩니다.
다음은 패키지 구성 형식입니다.

common 패키지
    _ST5_DATA/
    _ST5_VERSION.txt

extension 패키지
    _ST5_DATA/
    _ST5_VERSION.txt
    _ST5_EXE.*
    _ST5_ICON.*

starter5 구성 요소
    _ST5_DATA/
    _ST5_EXTENSION/
    _ST5_COMMON/
    _ST5_VERSION.txt
    _ST5_EXE.*
    _ST5_ICON.*
    _ST5_CONFIG.txt
    _ST5_SIGN.txt
