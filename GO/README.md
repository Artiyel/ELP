# Traitement d'image parallèle en Go

Ce projet implémente différents algorithmes de traitement d'image en exploitant la parallélisation offerte par le langage Go. L'objectif est d'optimiser les temps de calcul sur des images haute résolution en utilisant un système de Worker Pool.

---

## Fonctionnalité principale : Sobel par couleurs

Le script **`sobel_par_couleurs.go`** est le coeur de ce projet. Contrairement à un filtre Sobel classique qui traite l'image en noir et blanc, cette version :
* Applique l'algorithme de détection de contours directement sur les canaux **Rouge, Vert et Bleu** simultanément.
* Préserve les variations de couleurs dans les contours pour un rendu visuel plus riche.

---

## Autres outils de traitement

Bien que le filtre Sobel couleur soit l'aboutissement du projet, d'autres scripts ont été développés pour explorer le traitement d'image :

* **Sobel niveaux de gris (`sobel_par_gris.go`)** : Une approche plus traditionnelle qui convertit chaque pixel en luminance avant de calculer les contours.
* **Séparation de canaux (`separation_canaux_image.go`)** : Un outil permettant d'isoler les composantes RGB d'une image en trois fichiers distincts. Il traite l'image par blocs de lignes pour une meilleure gestion de la mémoire.

---

## Architecture : parallélisation et performance

Pour garantir une exécution rapide, tous les scripts reposent sur un modèle de concurrence robuste :

1.  **Worker Pool** : Le programme lance un nombre de "workers" (Goroutines) égal, par défaut, au nombre de cœurs logiques de votre processeur.
2.  **Gestion des jobs** : L'image est découpée en lignes (ou blocs de lignes) envoyées via un canal (`chan Job`).
3.  **Synchronisation** : Un `sync.WaitGroup` assure que toutes les parties de l'image sont traitées avant de lancer la sauvegarde.

---

## Mesure du temps (timer)

Le projet inclut une mesure précise du temps d'exécution pour comparer l'efficacité de la parallélisation. À chaque exécution, le programme affiche :
* **Temps de calcul** : Le temps réel passé par les workers à transformer l'image.
* **Temps de sauvegarde** : Le temps nécessaire pour encoder et écrire le fichier sur le disque.

*Note : La sauvegarde des images dans `separation_canaux_image.go` est également parallélisée pour réduire le temps total.*

---

## Utilisation

### Prérequis
* Go installé sur votre machine.
* Un dossier nommé `images/` contenant vos fichiers `.jpg`.

### Commandes
Vous pouvez exécuter les scripts directement avec la commande `go run`.

```bash
# Exécuter le traitement principal (Sobel Couleur) sur ma_photo.jpg, en utilisant 8 coeurs
go run sobel_par_couleurs.go -n 8 -f "images/ma_photo.jpg"

# Exécuter la séparation des canaux
go run separation_canaux_image.go -s=true