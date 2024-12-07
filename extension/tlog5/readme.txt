extension.tlog5 [win & linux] : 로그기록 관리도구.

<Auto Log>
자동으로 가장 마지막 로그 넘버 + 1에 해당하는 로그를 기록.

<Manual Log>
로그 넘버를 수동으로 지정하여 기록.

<Merge Log>
두 로그(main & sub)를 로그 넘버를 기준으로 병합. 결과는 main 로그 파일에 하나로 합쳐짐.

<Fork Log>
특정 로그 넘버보다 큰 로그만을 데이터 폴더에 new_branch.txt로 저장.

timestamp와 lognum, info로 이루어진 로그 라인이 모여 로그 기록 파일이 됩니다.
tlog5는 로그 기록(branch)을 관리하고 로그 추가 / 병합(merge) / 포크 등을 할 수 있습니다.
