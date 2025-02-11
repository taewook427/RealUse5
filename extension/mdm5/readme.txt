extension.mdm5 [win & linux] : 5세대 미니멀 데이터 매니저.

<지원 형식 gen1>
KENC : 간단한 파일 암호화 형식.

<지원 형식 gen2>
KENC : 간단한 파일 암호화 형식.

<지원 형식 gen3>
KZIP : 파일/폴더 패키징 형식.
KAES : 사용자용 파일 암호화 (all), KV용 데이터 암호화 (func)
KPIC : 파일을 픽셀 데이터 변조로 사진으로 패키징.
KV3st : 폴더 전체를 암호화/복호화 - 쉘 지원하지 않음. (G3KAES 사용)
ZIP release : 파일을 겉보기에 사진인 zip 파일로 패키징.

<지원 형식 gen4>
KENC : 간단한 파일 암호화 형식.
KAES : 사용자용 파일 암호화 (all), KV용 데이터 암호화 (func)
KV4st : 폴더 전체를 암호화/복호화 - 쉘 지원하지 않음. (G4KAES 사용)

<지원 형식 gen5>
File Div : 하나의 큰 파일을 여러 개의 작은 파일로 분할.
KSC : RealUse5의 데이터 직렬화 형식.
KZIP : 파일과 폴더를 하나의 파일로 패키징하는 형식.
KAES : 사용자용 파일 암호화 (all), KV용 데이터 암호화 (func)
KPKG : OS별로 배포할 수 있는 패키지 형식.
KPIC : 파일을 픽셀 데이터 변조로 사진으로 패키징.
kscript : kscript5 코드의 컴파일과 리버싱 지원.
KV4adv : RealUse5의 파일 기반 클러스터 - 쉘 지원.
KV5st : RealUse5의 청크 기반 클러스터 - 쉘 지원.

<KV4adv, KV5st Shell>
RealUse5 kvault의 클러스터는 기본적인 파일 입출력 기능을 쉘로 지원합니다.
쉘 명령어는 [명령] [옵션] [옵션]... 과 같이 공백으로 구분하여 사용합니다.

cd : 현재 디렉토리 이동 / 새로고침 #옵션 - 없음 : 새로고침, -1 : 상위 폴더로, N[숫자] : 해당 폴더로.
rm : 파일/폴더 삭제 #옵션 - N[숫자] : 해당 항목 삭제.
ren : 파일/폴더 이름바꾸기 #옵션 - N[숫자] S[새이름] : 해당 항목 이름을 S로 바꿈.
mv : 파일/폴더 이동 #옵션 - N[숫자]... : 해당 항목들을 이동시킴. #이동 모드 명령어 - -2 : 대상 폴더 선택, -1 : 상위 폴더로, N[숫자] : 해당 폴더로.
md : 새 폴더 생성. #옵션 - S[이름] : 새 폴더 S 생성.
touch : 새 파일 생성. #옵션 - S[이름] : 새 파일 S 생성.
im : 외부 파일/폴더 가져오기. #옵션 - P[경로] : P 경로의 파일/폴더를 클러스터로 가져옴.
ex : 파일/폴더 내보내기. #옵션 - 없음 : 현재 디렉토리 내보내기, N[숫자] : 해당 항목 내보내기.
view : 상태 보기. #옵션 - 없음 : 비동기 동작과 오류, log : 로그 데이터, debug : 디버그 정보, print : 모든 하위 파일과 폴더.
reset : 폴더 크기 보기 설정 / 비밀번호 재설정 #옵션 - t : 폴더 크기 보기, f : 폴더 크기 안 봄, acc : 비밀번호 재설정.
exit : 쉘 종료.
기타 : 도움말 출력.

일부 명령어는 다른 키워드로도 작동시킬 수 있습니다.
del -> rm, move -> mv, mkdir -> md, mkfile -> touch, import -> im, export -> ex
true/on -> t, false/off -> f, account/pw/pwkf -> acc

단일 프로그램 하나로 대부분의 데이터형식을 다룰 수 있게 해줍니다.
RealUse5 뿐만 아니라, 초창기 파이썬 코드들이 사용했던 형식(1세대 ~ 4세대)도 지원합니다.
Starter5 환경에서 동작하지만, 단일 실행파일로도 일부 기능을 제외하면 사용할 수 있습니다.
