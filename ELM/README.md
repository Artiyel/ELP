# ELM : Jeu "Guess It!"


## Concept du jeu
Le but est de deviner un mot pioché aléatoirement dans une liste de mots courants. 
* L'application charge une liste de mots depuis un fichier texte serveur.
* Un mot est choisi au hasard via un générateur d'index.
* Les définitions du mot sont récupérées en temps réel via une API externe pour aider l'utilisateur.
* L'utilisateur saisit sa proposition et l'interface valide la réponse dynamiquement.

---

## Architecture du code

Le projet est découpé en modules spécialisés pour garantir une maintenance facile et une séparation des responsabilités :

### 1. Gestion des données et API (`Fetch_def.elm`)
Ce module gère toute la communication avec l'API *Free Dictionary*.
* **Décodage JSON** : Utilise des décodeurs (`map2`, `map3`, `list`) pour transformer la réponse complexe de l'API en structures Elm exploitables (`Definition`, `Meaning`).
* **Requêtes HTTP** : Envoie des requêtes `GET` vers l'API d'après le mot choisi.

### 2. Logique de sélection (`Choose.elm`)
S'occupe de la manipulation du texte brut pour en extraire des mots :
* **Générateur** : Crée un index aléatoire basé sur le nombre d'espaces trouvés dans le texte source.
* **Extraction** : Découpe le texte (`String.slice`) pour isoler le mot correspondant à l'index généré.

### 3. Comparaison et Validation (`Compare.elm`)
Module utilitaire contenant la fonction `isTheSame`. Il permet une comparaison insensible à la casse (`String.toLower`) pour ne pas pénaliser l'utilisateur sur les majuscules.

### 4. Coeur de l'application (`Main.elm`)
Il assemble les pièces selon l'architecture Elm :
* **Model** : Stocke l'état du jeu (mot cible, saisie utilisateur, définitions, score).
* **Update** : Gère les messages comme le chargement du texte (`GotText`), la réception des définitions (`FetchMsg`) ou le changement de saisie (`Change`).
* **View** : Une interface utilisateur claire avec des retours visuels colorés (vert pour juste, rouge pour faux).

---

## Installation et lancement

1. Assurez-vous d'avoir installé **Elm**.
2. Entrez la commande suivante dans le terminal depuis le répertoire racine `ELM` :

   ```bash
   elm reactor
3. Une fois dans le navigateur, rendez vous dans le répertoire `src` et lancez `Main.elm`.
