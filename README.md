# API du Dictionnaire

Bienvenue dans l'API du dictionnaire ! Cette API permet d'ajouter, définir, supprimer et lister des mots avec leurs définitions.

## Choisissez un mode :
1. Console
2. API

## Endpoints de l'API :

- **/api/login** : Attend une requête HTTP de type POST avec les informations d'identification (username et password) dans le corps de la requête. Si les informations sont valides, elle renvoie un jeton d'authentification.
{"username": "nabil", "password":"10"}

- **/api/words/list** : Attend une requête HTTP de type GET. Nécessite un jeton d'authentification pour obtenir la liste des mots.

- **/api/words/add** : Attend une requête HTTP de type POST avec les données du mot et de sa définition dans le corps de la requête (Word, Definition). Nécessite un jeton d'authentification pour ajouter un nouveau mot.

- **/api/words/define/** : Attend une requête HTTP de type PUT avec le mot spécifié dans l'URL (define/mot) et la nouvelle définition dans le corps de la requête. Nécessite un jeton d'authentification pour définir ou mettre à jour la définition d'un mot existant.

- **/api/words/remove/** : Attend une requête HTTP de type DELETE avec le mot spécifié dans l'URL (remove/mot). Nécessite un jeton d'authentification pour supprimer un mot.

## Démarrage du Serveur

Pour démarrer le serveur, exécutez la commande suivante :

```bash
go run main.go [mode]
```
Choisissez le mode en remplaçant [mode] par 1 pour la console ou 2 pour l'API.

## Tester

Pour tester l'application, exécutez la commande suivante :
```bash
go tests -v ./tests/
```
