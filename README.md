## Rate Limiter em Go

### Objetivo:
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

### Descrição: 
O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. </br> O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

#### Endereço IP: 
O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
#### Token de Acesso: 
O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens.</br> 
O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.

### Requisitos:
O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web.</br>
O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.</br>
O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.</br>
As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo .env na pasta raiz.</br>
Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.</br>
O sistema deve responder adequadamente quando o limite é excedido:
Código HTTP: 429
Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame</br>
Todas as informações de "limiter" devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.</br>
Crie uma "strategy" que permita trocar facilmente o Redis por outro mecanismo de persistência.</br>
A lógica do limiter deve estar separada do middleware.</br>

#### Exemplos:
Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.</br>
Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.</br>
Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

### Dicas:
Teste seu rate limiter sob diferentes condições de carga para garantir que ele funcione conforme esperado em situações de alto tráfego.

### Entrega:
O código-fonte completo da implementação.</br>
Documentação explicando como o rate limiter funciona e como ele pode ser configurado.</br>
Testes automatizados demonstrando a eficácia e a robustez do rate limiter.</br>
Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.</br>
O servidor web deve responder na porta 8080.


## Configuração

As configurações estão no arquivo [config.yml](config.yaml), sendo a principais:

default_rate_limit: 10 - configuração da quantidade de tentativas para expirar</br>
default_time_blocked: 10 - configuração do tempo padrão para aguardar liberar as tentativas em segundos, no exemplo 10 segundos</br>

## Configuração de Keys

As configurações de keys estão no arquivo [token-list.json](token-list.json) </br>
Como é um arquivo json, as chaves utilizadas dentro são:

token - configuração da chave que deve ser utilizada no API_KEY </br>
Exemplo: 01936972-68a4-7253-bf19-91a63080cc39

expiresIn - configuração do tempo que vai aguardar em segundos para liberar quando for bloqueado </br>

qtdRequests - quantidades de requisições que devem ser feitas para bloquear a key

## Como executar e testar a aplicação
### Rodando os testes

#### Docker
rodar o comando docker exec -it api make test

#### Terminal
rodar o comando go test .\test\ -v

### Rodando a aplicação

#### Docker
Rodar o comando: docker compose up -d --build

#### Terminal
Executar o comando: docker compose up redis -d
Executar o comando: go mod download
Executar o comando: go run main.go

### Fazendo o consumo da API

#### Por IP:

Para consumir a API por ip é só realizar uma chamada via GET para URL http://localhost:8080.

#### Por API Key:

Para consumir a API usando um key, é só realizar uma chamada via GET para URL http://localhost:8080, passando no header uma chave através do atribudo API_KEY