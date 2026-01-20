module Test exposing (..)

addElemInList lenombre fois danslaliste =
    if fois == 0 then danslaliste
    else addElemInList lenombre (fois-1) (lenombre::danslaliste)

dupli liste =
    case liste of
    (x::xs) -> x::x:: dupli xs
    [] -> []

compress liste =
    case liste of
    [] -> []
    [x]->[x]
    a::b::xs->
        if a==b then compress (b::xs)
        else a :: compress (b::xs)

addNR lenombre fois danslaliste =
    let a = List.repeat fois lenombre in
    List.append a danslaliste

intoliste nombre = [nombre,nombre]

dupliNR liste = 
    List.concatMap intoliste liste 

compressHelper x partialRes = 
    case partialRes of
    [] -> [x]
    (y :: ys) -> 
        if x == y
            then partialRes
        else x :: partialRes 

compressNR liste = liste 

--a finir a la prochaine fois

--TD2

type Couleur = Rouge | Noir
type Point = DonnePoint Float Float
type Maybe a = Just a | Nothing
type Result error value = Ok value | Err error
type StackInt = Empty | Push Int StackInt

type CouleurCarte = Pic|Carreaux|Coeur|Trefle
type ValeurCarte = As|Roi|Dame|Valet|Numero Int
type Carte = Carte ValeurCarte CouleurCarte
--Carte As trefle
listeAs = [Carte As Trefle, Carte As Pic, Carte As Carreaux, Carte As Coeur]
type Arbrebin a  =  Vide|Node a (Arbrebin a) (Arbrebin a)
arbre3float = Node 1.5 (Node 2.2 (Vide) (Vide)) (Node 3.5 (Vide) (Vide))

hauteur arbre = 
    case arbre of 
    Vide -> 0
    Node _ fg fd -> 1 + max (hauteur fg) (hauteur fd)