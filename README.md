# Sistema de Leilão em Go

## 📋 Introdução

Este projeto é um sistema de leilão desenvolvido em Go (Golang) que permite a criação e gerenciamento de leilões online. O sistema oferece funcionalidades para criar leilões, fazer lances, consultar usuários e determinar vencedores automaticamente. 

### Principais Funcionalidades:
- **Criação de Leilões**: Cadastro de produtos para leilão com diferentes condições (novo, usado, recondicionado)
- **Sistema de Lances**: Permite que usuários façam lances em leilões ativos
- **Processamento em Lote**: Otimização de performance através de inserção de lances em lotes
- **Finalização Automática**: Leilões são automaticamente finalizados após um intervalo configurável
- **API RESTful**: Interface completa para integração com aplicações frontend

### Arquitetura:
- **Clean Architecture**: Separação clara entre camadas de domínio, casos de uso e infraestrutura
- **MongoDB**: Banco de dados NoSQL para persistência
- **Gin Framework**: Framework web rápido e minimalista para Go
- **Docker**: Containerização para facilitar deployment e desenvolvimento

---

## 🔧 Requisitos do Sistema

### Ambiente de Desenvolvimento Local:

#### Obrigatórios:
- **Go**: versão 1.20 ou superior
- **MongoDB**: versão 4.4 ou superior
- **Git**: para controle de versão

#### Opcionais (para desenvolvimento com Docker):
- **Docker**: versão 20.10 ou superior
- **Docker Compose**: versão 2.0 ou superior

### Verificação dos Requisitos:

```bash
# Verificar versão do Go
go version

# Verificar se o MongoDB está instalado
mongod --version

# Verificar Docker (se usando containerização)
docker --version
docker-compose --version
```

---

## 🚀 Configuração e Execução Local

### 1. Instalação de Dependências

```bash
# Clone o repositório
git clone git@github.com:lucasfeitozas/golang-leilao-routine.git
cd golang-leilao-routine

# Baixar dependências do Go
go mod download

# Verificar se todas as dependências foram instaladas
go mod verify
```

### 2. Configuração de Variáveis de Ambiente

Crie ou edite o arquivo `.env` em `cmd/auction/.env`:

```bash
# Configurações de Leilão
BATCH_INSERT_INTERVAL=20s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=20s

# Configurações do MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=admin
MONGODB_URL=mongodb://admin:admin@localhost:27017/auctions?authSource=admin
MONGODB_DB=auctions
```

### 3. Configuração do MongoDB Local

```bash
# Iniciar MongoDB (Ubuntu/Debian)
sudo systemctl start mongod

# Iniciar MongoDB (macOS com Homebrew)
brew services start mongodb-community

# Criar usuário admin no MongoDB
mongosh
use admin
db.createUser({
  user: "admin",
  pwd: "admin",
  roles: ["userAdminAnyDatabase", "dbAdminAnyDatabase", "readWriteAnyDatabase"]
})
```

### 4. Executar a Aplicação

```bash
# Executar a partir do diretório raiz
go run cmd/auction/main.go

# Ou compilar e executar
go build -o auction cmd/auction/main.go
./auction
```

A aplicação estará disponível em: `http://localhost:8080`

---

## 🐳 Guia Docker/Docker-Compose

### Pré-requisitos de Instalação do Docker

#### Ubuntu/Debian:
```bash
# Atualizar repositórios
sudo apt update

# Instalar Docker
sudo apt install docker.io docker-compose

# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER

# Reiniciar sessão ou executar
newgrp docker
```

#### macOS:
```bash
# Instalar Docker Desktop
# Baixar de: https://www.docker.com/products/docker-desktop

# Ou via Homebrew
brew install --cask docker
```

#### Windows:
```bash
# Instalar Docker Desktop
# Baixar de: https://www.docker.com/products/docker-desktop
```

### Configuração dos Containers

O projeto já inclui os arquivos de configuração necessários:

- **Dockerfile**: Define a imagem da aplicação Go
- **docker-compose.yml**: Orquestra os serviços (aplicação + MongoDB)

### Comandos para Construir e Iniciar os Serviços

```bash
# Construir e iniciar todos os serviços
docker-compose up --build

# Executar em background (modo detached)
docker-compose up -d --build

# Apenas iniciar (sem rebuild)
docker-compose up

# Parar todos os serviços
docker-compose down

# Parar e remover volumes (limpar dados)
docker-compose down -v

# Ver logs dos serviços
docker-compose logs

# Ver logs de um serviço específico
docker-compose logs app
docker-compose logs mongodb
```

### Comandos Úteis para Desenvolvimento

```bash
# Reconstruir apenas a aplicação
docker-compose build app

# Executar comandos dentro do container
docker-compose exec app sh

# Acessar MongoDB via container
docker-compose exec mongodb mongosh -u admin -p admin

# Ver status dos containers
docker-compose ps

# Reiniciar um serviço específico
docker-compose restart app
```

### Acesso à Aplicação

Após inicialização bem-sucedida:

- **API da Aplicação**: `http://localhost:8080`
- **MongoDB**: `localhost:27017`
- **Credenciais MongoDB**: admin/admin

---

## 📡 Informações Adicionais

### Portas Utilizadas

| Serviço | Porta | Descrição |
|---------|-------|-----------|
| API Go | 8080 | Servidor web principal |
| MongoDB | 27017 | Banco de dados |

### Endpoints Principais

#### Leilões (Auctions)
```bash
# Listar leilões
GET /auction?status=0&category=eletrônicos&productName=smartphone

# Buscar leilão por ID
GET /auction/{auctionId}

# Criar novo leilão
POST /auction
Content-Type: application/json
{
  "product_name": "iPhone 14",
  "category": "eletrônicos",
  "description": "iPhone 14 em excelente estado",
  "condition": 1
}

# Buscar lance vencedor
GET /auction/winner/{auctionId}
```

#### Lances (Bids)
```bash
# Criar lance
POST /bid
Content-Type: application/json
{
  "user_id": "uuid-do-usuario",
  "auction_id": "uuid-do-leilao",
  "amount": 1500.00
}

# Buscar lances por leilão
GET /bid/{auctionId}
```

#### Usuários (Users)
```bash
# Buscar usuário por ID
GET /user/{userId}
```

### Códigos de Status

- **0**: Leilão Ativo
- **1**: Leilão Finalizado

### Condições de Produto

- **1**: Novo
- **2**: Usado  
- **3**: Recondicionado

### Solução de Problemas Comuns

#### Problema: Erro de conexão com MongoDB
```bash
# Verificar se MongoDB está rodando
docker-compose ps

# Verificar logs do MongoDB
docker-compose logs mongodb

# Reiniciar serviço MongoDB
docker-compose restart mongodb
```

#### Problema: Porta 8080 já está em uso
```bash
# Verificar processos usando a porta
lsof -i :8080

# Matar processo específico
kill -9 <PID>

# Ou alterar porta no docker-compose.yml
ports:
  - "8081:8080"  # Usar porta 8081 externamente
```

#### Problema: Erro ao compilar aplicação Go
```bash
# Limpar cache do Go
go clean -modcache

# Baixar dependências novamente
go mod download

# Verificar versão do Go
go version  # Deve ser 1.20+
```

#### Problema: Variáveis de ambiente não carregadas
```bash
# Verificar se arquivo .env existe
ls -la cmd/auction/.env

# Verificar conteúdo do arquivo
cat cmd/auction/.env

# Recriar arquivo se necessário
cp cmd/auction/.env.example cmd/auction/.env
```

#### Problema: Erro de permissão Docker (Linux)
```bash
# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER

# Reiniciar sessão ou executar
newgrp docker

# Verificar se funcionou
docker ps
```

### Logs e Debugging

```bash
# Ver logs da aplicação
docker-compose logs -f app

# Ver logs com timestamp
docker-compose logs -t app

# Executar aplicação em modo debug local
go run cmd/auction/main.go

# Verificar conectividade com MongoDB
docker-compose exec app ping mongodb
```

### Performance e Monitoramento

O sistema inclui otimizações de performance:

- **Inserção em Lote**: Lances são processados em lotes para melhor performance
- **Timers Assíncronos**: Finalização automática de leilões sem bloquear requests
- **Conexão Pooling**: MongoDB driver otimizado para múltiplas conexões

Para monitorar performance:
```bash
# Verificar uso de recursos
docker stats

# Monitorar logs em tempo real
docker-compose logs -f
```

---

## 🤝 Contribuição

Para contribuir com o projeto:

1. Faça fork do repositório
2. Crie uma branch para sua feature (`git checkout -b feature/nova-funcionalidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova funcionalidade'`)
4. Push para a branch (`git push origin feature/nova-funcionalidade`)
5. Abra um Pull Request

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.