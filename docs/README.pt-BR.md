# GonGo-Simple-Auth

### Visão Geral do Projeto

O projeto foi concebido para ser um *boilerplate* de autenticação direto e pronto para uso. Ele implementa registro de usuários, login com emissão de JWT e uma rota protegida `/me` para verificar a validade da sessão.

### Funcionalidades Principais

* **Gestão de Usuários**: Gerencia registro e login utilizando MongoDB para persistência.
* **Segurança**: Utiliza `bcrypt` para hashing de senhas e `HMAC-SHA256` para geração de tokens JWT.
* **Proteção por Middleware**: Inclui um middleware `Guard` que valida tokens Bearer em rotas restritas.
* **Configuração de Ambiente**: Gerenciada via arquivos `.env` com um unmarshaler customizado.
* **Ferramentas de Desenvolvimento**: Pré-configurado com `Air` para live reloading (atualização em tempo real).

### Stack Técnica

* **Linguagem**: Go 1.25.5.
* **Framework Web**: Gin Gonic.
* **Banco de Dados**: MongoDB Driver.
* **Autenticação**: JWT-Go (v5).

### Como Executar

1. Configure seu arquivo `.env` com as variáveis `GONGO_MONGO_URI`, `GONGO_PORT` e `GONGO_JWT_SECRET`.
2. Instale as dependências: `go mod tidy`.
3. Execute com Air: `air` (ou `go run cmd/api/main.go`).
