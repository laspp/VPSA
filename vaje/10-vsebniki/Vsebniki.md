# Vsebniki

Pogosto se zgodi, da na sistemu nimamo nameščenih vseh zahtevanih programov in knjižnic. V določenih okoljih, na primer na računski gruči Arnes, uporabniki pogosto nimajo ustreznih dovoljenj za nameščanje dodatnega programja. Prav tako se včasih želimo izogniti posegom v sistemsko okolje in nameščanju dodatnih paketov, ki bi lahko vplivali na obstoječo konfiguracijo. V takšnih primerih nam pridejo zelo prav vsebniki.

## Vsebniki Apptainer

Na superračunalniških gručah SLING se za upravljanje z vsebniki uporablja ogrodje [Appatiner](https://apptainer.org/). Poglejmo si, kako z njegovo pomočjo na gruči Arnes pripravimo vsebnik, znotraj katerega bomo lahko prevajali in zaganjali programe v Go, ki uporabljajo grpc.

Apptainer je na gruči Arnes že nameščen. Preverimo, verzijo trenutne namestitve:

```bash
$ apptainer version
1.4.5-2.el9
```

Apptainer pozna:

- **Nespremenljive vsebnike** (angl. immutable) imajo končnico `.sif` in so namenjeni uporabi v produkciji. Ko jih enkrat zgradimo, jih ne moremo več spreminjati.
- **Spremenljive vsebnike** (angl. mutable) so namenjeni predvsem razvoju in testiranju. Tak vsebnik je na gostitelju predstavljen z mapo, v kateri so shranjene vse datoteke vsebnika. Po zaključeni gradnji in testiranju spremenljiv vsebnik običajno pretvorimo v nespremenljivega.

### Gradnja spremenljivega vsebnika


Začnimo z gradnjo spremenljivega vsebnika na osnovi operacijskega sistema Ubuntu. Spremenljive vsebnike bomo označevali s končnico `.sb` (angl. Sandbox).

```bash
$ apptainer build --sandbox go-grpc.sb docker://ubuntu:24.04
```

Če želimo vsebnik spreminjati, moramo imeti skrbniške pravice. Na večuporabniškem sistemu uporabimo stikalo `--fakeroot`, ki nam omogoča delo s skrbniškimi pravicami brez pravega dostopa `sudo`. Prav tako moramo označiti, da želimo vsebnik odpreti za pisanje s stikalom `--writable`. Zaženimo ukazno lupino znotraj vsebnika:

```bash
$ apptainer shell --fakeroot --writable go-grpc.sb
```

Najprej posodobimo paketni sistem in namestimo potrebne pakete:

```bash
Apptainer> apt update
Apptainer> apt upgrade -y
Apptainer> apt install -y wget build-essential unzip
```

Namestimo najnovejšo različico programskega jezika Go:

```bash
Apptainer> export SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
Apptainer> cd /usr/local
Apptainer> wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
Apptainer> tar -xzf go1.25.5.linux-amd64.tar.gz
Apptainer> rm go1.25.5.linux-amd64.tar.gz
```
Nastavimo okoljske spremenljivke za trenutno sejo:

```bash
Apptainer> export PATH=$PATH:/usr/local/go/bin
Apptainer> export GOPATH=/root/go
Apptainer> export PATH="$PATH:$GOPATH/bin"
```

Preverimo namestitev:

```bash
Apptainer> go version
go version go1.25.5 linux/amd64
```

Namestimo Protocol Buffers compiler (protoc):

```bash
cd /usr/local
wget https://github.com/protocolbuffers/protobuf/releases/download/v33.2/protoc-33.2-linux-x86_64.zip
unzip protoc-33.2-linux-x86_64.zip
rm protoc-33.2-linux-x86_64.zip
```

Namestimo module Go za gRPC in Protocol Buffers:

```bash
Apptainer> go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
Apptainer> go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

```

Preverimo namestitev:

```bash
Apptainer> protoc --version
libprotoc 33.2

Apptainer> which protoc-gen-go
/root/go/bin/protoc-gen-go

Apptainer> which protoc-gen-go-grpc
/root/go/bin/protoc-gen-go-grpc
```

Ko smo končali z namestitvami, izstopimo iz vsebnika:

```bash
Apptainer> exit
```
Okoljske spremenljivke smo v vsebniku nastavili samo za trenutno sejo. Če želimo, da se bodo ustrezno nastavile tudi ob naslednjem zagonu vsebnika jih dodamo v datoteko `go-grpc.sb/.singularity.d/env/90-environment.sh`:

```bash
$ cat >> go-grpc.sb/.singularity.d/env/90-environment.sh << 'EOF'
export SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/go
export PATH=$PATH:$GOPATH/bin
EOF
```

Preverimo, ali lahko v vsebniku uporabljamo Go:

```bash
$ apptainer exec go-grpc.sb go version
go version go1.25.5 linux/amd64
```

Ko smo s spremenljivim vsebnikom zadovoljni, ga lahko pretvorimo v nespremenljiv vsebnik `.sif`:

```bash
$ apptainer build go-grpc.sif go-grpc.sb
```

Po uspešni pretvorbi lahko spremenljiv vsebnik izbrišemo:

```bash
$ rm -rf go-grpc.sb
```

Preizkusimo nespremenljiv vsebnik:

```bash
$ apptainer exec go-grpc.sif go version
go version go1.25.5 linux/amd64
```

Poskusimo pognati strežnik grpc znotraj vsebnika. Uporabimo lahko kar predlogo iz [vaje 8](./koda/predloga/). Najprej poskrbimo, da se mapa s kodo (`predloga`) nahaja v isti mapi kot naš vsebnik, potem pa inicializiramo modul, namestimo ustrezne pakete in poženemo strežnik s spodnjim ukazom:
```bash
srun --reservation=fri apptainer exec go-grpc.sif bash -c "cd ./predloga && go mod init api && go mod tidy && cd ./grpc && go run ."
```

Če smo uspešni bi se moral postaviti strežnik grpc na računskem vozlišču:
```
gRPC server listening at wn101.arnes.si:9876
```

## Domača naloga 8

V prejšnjem delu smo vsebnik zgradili interaktivno - vstopili smo vanj in ročno izvedli vse potrebne ukaze. V praksi je bolj primerno uporabljati recepte (angl. definition files), ki omogočajo avtomatizirano in ponovljivo gradnjo vsebnikov. Pripravite recept `go-grpc.def`, ki bo avtomatično zgradil vsebnik z nameščenim Go in gRPC. Recept naj izvede enake korake kot smo jih izvedli interaktivno. Rešitev (recept) oddajte preko [spletne učilnice](https://ucilnica.fri.uni-lj.si/mod/assign/view.php?id=61268).

### Struktura recepta

Recept za Apptainer vsebnik ima naslednje glavne sekcije:

```
Bootstrap: docker
From: ubuntu:24.04

%post
    # Ukazi, ki se izvedejo med gradnjo vsebnika
    # Tukaj namestite vse potrebne pakete in programsko opremo

%environment
    # Definicije okoljskih spremenljivk
    # Te spremenljivke bodo na voljo ob zagonu vsebnika

%runscript
    # Ukaz, ki se izvede, ko zaženemo vsebnik z "apptainer run"

%labels
    # Metapodatki o vsebniku (avtor, različica, opis)

%help
    # Besedilo, ki se prikaže, ko uporabnik požene "apptainer run-help"
```

### Navodila za izvedbo

1. Ustvarite datoteko `go-grpc.def` z ustrezno vsebino recepta.

2. V sekciji `%post` zapišite vse ukaze za:
   - Posodobitev paketnega sistema
   - Namestitev osnovnih orodij (wget, build-essential, unzip)
   - Namestitev Go (prenesite in razpakirajte v `/usr/local`)
   - Namestitev prevajalnika Protocol Buffers (prenesite in razpakirajte v `/usr/local`)
   - Namestitev modulov za gRPC (`protoc-gen-go` in `protoc-gen-go-grpc`)

3. V sekciji `%environment` definirajte okoljske spremenljivke:
   - `PATH` naj vključuje pot do Go (`/usr/local/go/bin`)
   - `GOPATH` naj kaže na `/root/go`
   - `PATH` naj vključuje pot do Go binark (`$GOPATH/bin`)
   - `SSL_CERT_FILE` naj kaže na `/etc/ssl/certs/ca-certificates.crt`

4. V sekciji `%runscript` dodajte ukaz, ki prikaže različico Go.

5. V sekciji `%labels` dodajte metapodatke:
   - Avtor
   - Različica
   - Opis vsebnika

6. V sekciji `%help` dodajte kratko navodilo za uporabo vsebnika.

7. Ko pripravite recept, zgradite vsebnik z ukazom:

```bash
$ apptainer build go-grpc.sif go-grpc.def
```

**Rok za oddajo 18. 1. 2026.**