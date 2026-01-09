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

  ```console
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

- primerjava izpisov potrdi, nekaj pomembnih lastnosti vsebnikov Apptainer:

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
