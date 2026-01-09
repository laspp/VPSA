# Delo z vsebniki Apptainer

- pri pripravi vsebnikov lahko uporabimo naslednje možnosti:

  - lahko jih potegnemo iz podprtih repozitorijev, kot so [Docker Hub](https://hub.docker.com/), [Syngularity Cloud Library](https://cloud.sylabs.io/library) in [Nvida NGC](https://ngc.nvidia.com/)
  - gradimo jih interaktivno v peskovniku ali interaktivno popravljamo obstoječe slike,
  - gradimo jih neinteraktivno iz v naprej pripravljenih receptov

## Uporabimo obstoječi vsebnik

- delujoči vsebnik dobimo najhitreje tako, da ga potegnemo iz repozitorija
  - `docker://` za vsebnike v repozitoriju [Docker Hub](https://hub.docker.com/)
  - `library://` za vsebnike v repozitoriju [Singularity Cloud Library](https://cloud.sylabs.io/library)
- uporabimo zabaven vsebnik `lolcow`, ki ga najdemo na repozitoriju [Docker Hub](https://hub.docker.com/r/sylabsio/lolcow)

  ```bash
    $ apptainer pull docker://sylabsio/lolcow:latest
  ```

- po končanem prenosu imamo v mapi datoteko `lolcow_latest.sif`; vsebniki Apptainer/Singularity so označeni s končnico `sif` (*angl.* Singularity Image File)

### Interaktivno delo v vsebniku

- preden zaženemo vsebnik, poglejmo nekaj nastavitev našega računalnika (gostitelja):

  ```bash
    $ hostname
    hpc-login1.arnes.si
    $ whoami
    sling001
    $ cat /etc/os-release
    NAME="AlmaLinux"
    VERSION="8.10 (Cerulean Leopard)"
    ...
    $ echo $PATH
    /d/hpc/home/sling001/.local/bin:/d/hpc/home/sling001/bin:/usr/local/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/opt/puppetlabs/bin
    $ ls 
    lolcow_latest.sif
    $ ls /usr/games
  ```

- za interaktivno delo z vsebnikom uporabimo ukaz `shell`, ki bo vsebnik zagnal v ukazni vrstici, na podoben način, kot pri povezovanju na oddaljeni sistem preko protokola SSH

  ```bash
  $ apptainer shell lolcow_latest.sif
  ```

- vpišimo še enkrat vse prejšnje ukaze in še nekaj dodatnih

  ```bash
  Apptainer> hostname
  hpc-login1.arnes.si
  Apptainer> whoami
  sling001
  Apptainer> cat /etc/os-release    
  NAME="Ubuntu"
  VERSION="20.04.2 LTS (Focal Fossa)"
  ...
  Apptainer> echo $PATH
  /usr/games:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
  Apptainer> ls 
  lolcow_latest.sif
  Apptainer> ls /usr/games
  cowsay  cowthink  lolcat
  Apptainer> cowsay "Hello!"
      --------
     < Hello! >
      --------
  ...
  Apptainer> cp /etc/os-release ./lolcow-os.info
  Apptainer> exit
  ```

- primerjava izpisov potrdi nekaj pomembnih lastnosti vsebnikov Apptainer:

  - vsebnik prevzame množico nastavitev od gostitelja (hostname)
  - uporabniško ime v vsebniku je enako uporabniškemu imenu na gostitelju; enako je s pravicami uporabnika
  - vsebnik zamenja uporabniški prostor (Ubuntu namesto AlmaLinux)
  - okoljska spremenljivka `PATH` se v vsebniku spremeni; vsebnik privzeto vidi mape `$HOME`, `$PWD`, `/sys`, `/proc`, `/tmp` in še nekatere, mape `/d/hpc/home/sling001/.local/bin`, ki ni na teh lokacijah, v vsebniku ne vidimo; vsebnik je spremenljivko `PATH` dopolnil z mapo `/usr/games`
  - datoteke v omenjenih mapah, izpisali smo datoteki v trenutni mapi (`$PWD`), vidimo tudi v vsebniku
  - mapa `/usr/games` na gostitelju je prazna, v vsebniku pa so v njej štirje programi
  - v vsebniku lahko zaganjamo ukaze operacijskega sistema in druge programe (`ls`, `cowsay`)
  - iz vsebnika lahko tudi pišemo (kopiramo) v datotečni sistem gostitelja (ukaz `cp`).

- z ukazom `exit` zaključimo interaktivno delo in se vrnemo v ukazno vrstico gostitelja.

### Izvajanje ukazov v vsebniku iz gostitelja

- z ukazom `exec` lahko izvajamo ukaze v vsebniku
- ogrodje Apptainer najprej zažene vsebnik, v našem primeru `lolcow_latest.sif`, nato pa v vsebniku zažene ukaz (`ls`) ali poljuben program, ki je nameščen v njem

  ```bash
  $ apptainer exec lolcow_latest.sif ls /usr/games/
  cowsay  cowthink  fortune  lolcat
  $ apptainer exec lolcow_latest.sif cowthink "Hmmmm..."
   ----------
  ( Hmmmm... )
   ----------
   ...
  $ apptainer exec lolcow_latest.sif lolcat lolcow-os.info
  ```

### Zaganjanje vsebnika

- pri pripravi vsebnika lahko v zagonsko skripto (*angl.* runscript) vpišemo ukaze, ki se izvedejo ob zagonu
- vsebnik zaženemo z ukazom `run` ali pa preprosto vpišemo njegovo ime, saj je izvršljiva datoteka

  ```bash
  $ apptainer run lolcow_latest.sif
   _______________________________
  < Mon Jun 23 11:46:53 CEST 2025 >
   -------------------------------
          \   ^__^
           \  (oo)\_______
              (__)\       )\/\
                  ||----w |
                  ||     ||     
  $ ./lolcow_latest.sif
   _______________________________
  < Mon Jun 23 11:48:15 CEST 2025 >
   -------------------------------
  ...
  $ apptainer run docker://sylabsio/lolcow:latest
   _______________________________
  < Mon Jun 23 11:48:56 CEST 2025 >
   -------------------------------
  ...
  ```

- v zadnjem primeru zaženemo vsebnik kar neposredno iz repozitorija: ogrodje Apptainer v tem primeru pripravi začasno sliko, ki jo po zaključku izvajanja vsebnika zbriše

## Interaktivna gradnja vsebnika

- zgradili bomo vsebnik zasnovan na operacijskem sistemu Linux **Ubuntu** in vanj namestili programček `fortune`.

  ```bash
  $ apptainer build modrec.sif docker://ubuntu
  ```

- vsebnik zaženimo v interaktivnem načinu

  ```bash
  $ apptainer shell modrec.sif
  ```

- poglejmo verzijo operacijskega sistema

  ```bash
  Apptainer> cat /etc/os-release
  NAME="Ubuntu"
  VERSION_ID="24.04"
  VERSION="24.04.2 LTS (Noble Numbat)"
  ...
  ```

### Nespremenljivi in spremenljivi vsebniki

- ogrodje Apptainer pozna nespremenljive (*angl.* immutable) in spremenljive (*angl.* mutable) vsebnike
- nespremenljive vsebnike spoznamo po datoteki s končnico `sif`; običajno delamo z njimi
- spremenljive vsebnike potrebujemo predvsem med razvojem (gradnjo)
- spremenljiv vsebnik zahtevamo s stikalom `--sandbox`; tak vsebnik je na gostitelju predstavljen z mapo, v kateri so shranjene vse datoteke vsebnika; označili ga bomo s končnico `sb` (*angl.* SandBox)
- po zaključeni gradnji spremenljiv vsebnik pretvorimo v nespremenljivega

- če želimo vsebnik spreminjati, moramo imeti skrbniške pravice - na lastnem računalniku uporabimo ukaz `sudo` (*angl.* super user do) ali stikalo `--fakeroot`, na večuporabniškem sistemu pa stikalo `--fakeroot`; dodati moramo tudi stikalo ``--writable`

  ```bash
  $ apptainer build --sandbox modrec.sb docker://ubuntu
  INFO:    Starting build...
  ...
  INFO:    Creating sandbox directory...
  INFO:    Build complete: modrec.sb

  $ ls modrec.sb
  bin   dev          etc   lib    lib64   media  opt   root  sbin         srv  tmp  var
  boot  environment  home  lib32  libx32  mnt    proc  run   singularity  sys  usr

  $ apptainer shell --fakeroot --writable modrec.sb

  Apptainer> apt update
  ...
  1 package can be upgraded. Run 'apt list --upgradable' to see it.

  Apptainer> apt install fortune
  ...
  The following NEW packages will be installed:
  fortune-mod fortunes-min librecode0
  0 upgraded, 3 newly installed, 0 to remove and 1 not upgraded.
  Need to get 762 kB of archives.
  After this operation, 2154 kB of additional disk space will be used.
  ...

  Apptainer> /usr/games/fortune
  It's always darkest just before it gets pitch black.

  Apptainer> exit

  $ apptainer exec modrec.sb /usr/games/fortune
  Never eat more than you can lift.

  $ apptainer run modrec.sb

  Apptainer> exit
  ```

- ukaz `ls -l modrec.sb` nam razkrije, da je spremenljiv vsebnik dejansko mapa z množico datotek
- v osnovni mapi vsebnika najdemo datoteko `singularity` - gre za skripto bash, ki se izvede ob zagonu vsebnika
- vsebino datoteke `singularity` zamenjajmo z

  ```bash linenums="1"
  #!/bin/sh
  exec /usr/games/fortune "$@"
  ```

- preverimo delovanje

  ```bash
  $ apptainer run modrec.sb
  There has been an alarming increase in the number of things you know nothing about.
  ```

- ko smo z vsebnikom zadovoljni, ga lahko pretvorimo v nespremenljivega

  ```bash
  $ apptainer build modrec.sif modrec.sb
  Build target 'modrec.sif' already exists and will be deleted during the build process. 
  Do you want to continue? [N/y]y
  INFO:    Starting build...
  INFO:    Creating SIF file...
  INFO:    Build complete: modrec.sif

  $ ./modrec.sif
  Computers are not intelligent. They only think they are.
  ```

## Recepti za vsebnike

- ko za gradnjo vsebnika pripravimo recept (*angl.* recipe / definition file / def file), zagotovimo, da bo vsakdo lahko vsebnik zgradil na enak način kot mi
- recept je sestavljen iz dveh delov:
  - iz glave, v kateri navedemo operacijski sistem, na katerem želimo vsebnik zgraditi (distribucijo in verzijo)
  - iz opcijskih razdelkov, v katere pišemo navodila za gradnjo in zaganjanje vsebnika

### Glava

- v glavi z zapisoma `Bootstrap: docker` in `From: ubuntu:latest` zahtevamo gradnjo vsebnika iz Linux distribucije `ubuntu`, najnovejše verzije, ki se nahaja v repozitoriju [Docker Hub](https://hub.docker.com/)

- nekaj pogostih nastavitev za `Bootstrap`:

  - `library`: če izhajamo iz vsebnika v repozitoriju  [Singularity Cloud Library](https://cloud.sylabs.io/library)
  - `docker`: če izhajamo iz vsebnika v repozitoriju [Docker Hub](https://hub.docker.com/)
  - `localimage`: če vsebnik gradimo iz že obstoječe lokalne slike vsebnika

- Za nastavitev `From` vpišemo v prvih dveh primerih oznako vsebnika v repozitoriju, v zadnjem primeru pa ime vsebnika na lokalnem računalniku

### Razdelki

- `%setup`: pišemo ukaze, ki jih bo ogrodje Singularity izvedlo na gostitelju potem, ko bo nameščen operacijski sistem v vsebniku
- `%files`: navedemo datoteke, ki jih želimo iz gostitelja prenesti v vsebnik; podatke za vsako datoteko pišemo v svojo vrstico, pri čemer najprej navedemo ime datoteke na gostitelju, nato pa pot oziroma ime datoteke v vsebniku
- `%post`: v tem delu nameščamo programsko opremo v vsebniku in s spleta prenašamo datoteke; med nameščanjem programi velikokrat nastavijo tudi okoljske spremenljivke; okoljske spremenljivke, navedene v razdelku `%post` veljajo samo v času nameščanja
- `%environment`: sem vpišemo nastavitve okoljskih spremenljivk, ki se morajo izvesti ob zagonu vsebnika; z nastavitvijo spremenljivke `LC_ALL=C` povozimo nastavljeno lokalizacijo z enostavno lokalizacijo `C` (*angl.* Computers) z ASCII kodno tabelo in ameriško angleščino. 
- `%runscript`: ukazi, napisani v tem razdelku, se bodo izvedli ob zagonu vsebnika z ukazom `apptainer run`
- `%test`: v razdelek `%test` napišemo ukaze, s katerimi preverimo pravilnost delovanja vsebnika; ukazi v razdelku se zaženejo po končani gradnji vsebnika ali ob klicu ukaza `apptainer test`
- `%labels`: osnovni podatki o vsebniku v obliki oznaka - vrednost; osnovne podatke o vsebniku izpišemo z ukazom `apptainer inspect`
- `%help`: sem napišemo navodila za uporabo vsebnika; uporabnik jih izpiše z ukazom `apptainer run-help`.

### Gradnja vsebnika

- primera:
  - recept za gradnjo vsebnika <a href="recepti/modrec.def" download>`modrec.def`</a>
  - recept za gradnjo vsebnika <a href="recepti/firbec.def" download>`firbec.def`</a>
- vsebnik iz recepta zgradimo s pomočjo ukaza `build` s skrbniškimi pravicami
- recept, po katerem je bil zgrajen vsebnik, lahko pridobimo iz vsebnika s pomočjo ukaza `apptainer inspect -d`

### Povezovanje z datotečnim sistemom gostitelja

- ogrodje Apptainer, za razliko od mnogih drugih ogrodij za gradnjo vsebnikov, poskrbi, da v vsebniku vidimo datotečni sistem gostitelja
- vsebnik privzeto vidi mape: `$HOME` (glavna mapa uporabnika), `$PWD` (trenutna mapa) , `/tmp` , `/proc` , `/sys` in `/dev`
- mape, katerih vsebino bi še radi videli v vsebniku, povežemo z vsebnikom ob klicu ukaza `apptainer exec` ali `apptainer run`, tako da navedemo stikalo `--bind` z argumentom `mapa-gostitelja:mapa-vsebnika`; če želimo povezati več map, pare `mapa-gostitelja:mapa-vsebnika` ločimo z vejicami.

- program `fortune` izpisuje modrosti, ki jih najde v vsebniku v mapi `/usr/share/games/fortunes`
- kadar ob klicu programa `fortune` dodamo argument z imenom datoteke z modrostmi, bo izbral modrost iz te datoteke
- program lahko izpisuje tudi naše modrosti, ki niso shranjene v vsebniku

  - pripraviti moramo dve datoteki: besedilno datoteko brez končnice, v kateri so modrosti ločene z znakom `%`, in za hitrejši dostop do modrosti še datoteko s končnico `dat`; slednjo ustvarimo iz besedilne datoteke z ukazom `strfile`, ki ga najdemo v našem vsebniku
  - predpostavimo, da sta datoteki nameščeni v mapi [modrosti](./recepti/modrosti/), ki je podmapa mape z vsebnikom

    ```bash
    $ apptainer run --bind ./modrosti:/usr/share/games/fortunes modrec.sif pregovori
    Sreča je zaveznica pogumnih.

    $ apptainer run --bind ./modrosti:/usr/share/games/fortunes modrec.sif riddles
    No fortunes found
    ```

## Dobre prakse pri gradnji vsebnikov

- programske pakete, programe, podatkovne datoteke in podobno nameščamo v mape operacijskega sistema in ne v mape, kot so na primer `/home`, `/tmp`
- vsebnik vedno ustrezno dokumentiramo, saj imajo uporabniki imajo radi vsebnike, ki jih znajo hitro uporabljati
- pomoč za uporabo lahko vključimo v recepte v `%runscript` ali `%help`
- v vsebnik vključimo kakšen enostaven testni primer iz katerega je razvidno tudi, kako ga uporabljati
- poskrbimo za pravilne nastavitve pravic dostopa do datotek; lastnik datotek v vsebniku mora biti sistemski račun
- poskrbimo, da datoteke z občutljivo vsebino, na primer `/etc/passwd`, `/etc/group`, `/etc/shadow` ne vsebujejo gesel
- produkcijske vsebnike vedno gradimo iz receptov in tako zagotovimo ponovljivost gradnje in delovanja
