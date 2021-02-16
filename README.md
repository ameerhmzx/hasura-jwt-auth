# Hasura JWT Authentication
Golang authentication server sample for hasura.

## Requirements
Docker, docker-compose & hasura-cli installed

## Instructions
- update `.env` & `hasura/config.yaml` file
- remove `db/.gitignore` as postgres container requires empty directory
- execute `docker-compose up`
- apply migrations `hasura migration apply`
- apply seeds `hasura seeds apply`
- apply metadata `hasura metadata apply`

After these steps you can register and login user using `register` and `login` graphql mutation respectively
Access Token can also be refreshed using `refreshToken` mutation.
Expiry time for both access token and refresh token can be modified in `.env` file
