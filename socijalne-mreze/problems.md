# Problems

## Zadaci za samostalnu vežbu

### 1. Propagacija informacija nakon uklanjanja centralnih čvorova

Za neku socijalnu mrežu koja je povezana, odnosno ima tačno jednu komponentu
povezanosti, kažemo da ima funkcionalnost propagacije informacija ako
informacija koju odašilje neki čvor može stići do proizvoljnog čvora u mreži,
pod uslovom da svaki čvor po prijemu informacije istu prosledi svim svojim
susedima.

Ispitati da li socijalne mreže iz `demo2` gube sposobnost propagacije
informacija, odnosno da li ostaju povezane, ako se iz mreže ukloni:

- čvor koji ima najveći betweenness centrality,
- čvor koji ima najveći closeness centrality,
- čvor koji ima najveći eigenvector centrality.

Pre toga verifikovati da svaka od mreža ima tačno jednu komponentu povezanosti.

Za uklanjanje čvora iz grafa koristiti metodu `RemoveNode(node_id)`, koja se
može pozvati nad proizvoljnim `graph.Graph` objektom. Pre uklanjanja čvora
napraviti kopiju grafa, jer će se uklanjati više čvorova po različitim
kriterijumima. Kopiju grafa je moguće napraviti koristeći funkciju
`CopyGraph(graph)` iz `demo8`.

### 2. Fazni prelaz u Erdos-Renyi modelu

Erdos-Renyi model slučajnog grafa u limitu, kada broj čvorova raste u
beskonačno, poseduje jedan fazni prelaz:

- ispod neke kritične verovatnoće povezanosti graf skoro sigurno nema
  gigantsku komponentu povezanosti, odnosno komponentu koja obuhvata većinu
  čvorova u mreži;
- iznad neke kritične verovatnoće povezanosti graf skoro sigurno ima
  gigantsku komponentu povezanosti.

Kritična verovatnoća povezanosti je ona koja rezultuje prosečnim stepenom čvora
grafa koji je jednak `1`.

Eksperimentalno verifikovati ovaj teorijski rezultat na sledeći način:

1. Generisati slučajne grafove neke fiksirane veličine, na primer `100`
   čvorova.
2. Varirati verovatnoću povezanosti od `0` do `1` sa korakom `0.001`.
3. Za svaki generisani graf odrediti veličinu najveće komponente povezanosti.
4. Prestati sa generisanjem random grafova onog trenutka kada generisani random
   graf ima tačno jednu komponentu povezanosti.

### 3. Small-world osobina i distribucija distanci

Verifikovati small-world osobinu socijalnih mreža iz `demo2`, odnosno formirati
i prikazati distribuciju distanci.

Distribucija distanci kazuje koliko parova čvorova se nalazi na distanci `1`,
na distanci `2`, na distanci `3` i tako redom do maksimalne distance. Broj
parova čvorova na nekoj distanci normalizovati ukupnim brojem parova čvorova
kako bi se dobila funkcija verovatnoće.

Za računanje distance za sve parove čvorova može se koristiti funkcija
`path.DijkstraAllPaths`.

### 4. Jake zajednice

Za neku zajednicu kažemo da je jaka ako za svaki čvor u zajednici važi da ima
više veza sa čvorovima iz njegove zajednice nego sa čvorovima koji nisu u
njegovoj zajednici.

Za socijalne mreže iz `demo2` i razne algoritme za detekciju zajednica
ustanoviti koje detektovane zajednice su jake, a koje nisu.

### 5. Predikcija inter-community linkova

Za neki link kažemo da je inter-community ili inter-klaster link ako je to link
koji spaja čvorove iz različitih zajednica.

Napraviti skriptu koja za socijalne mreže iz `demo2` radi nesuperviziranu
predikciju inter-community linkova koristeći:

- Adamic-Adar metriku sličnosti čvorova,
- Louvain algoritam za detekciju zajednica.
