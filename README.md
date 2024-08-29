## Motivação

Este é um Laboratório da Pós-Graduação da Fullcycle: Go Expert, Desenvolvimento Avançado em Golang

## Laboratório
**Objetivo:**
Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

Clone o seguinte repositório: clique para acessar o repositório.

Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

**Você deverá desenvolver:**

- Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
- Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
- Um teste para validar se o fechamento está acontecendo de forma automatizada;

**Dicas:**

- Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
- Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso:
- Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;
- Para mais informações de como funciona uma goroutine, clique aqui e acesse nosso módulo de Multithreading no curso Go Expert;
 
**Entrega:**

- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.

## Como Executar

### Requisitos

- Makefile
- Docker
- Docker Compose
- Curl

### Rodar o sistema em ambiente de DESENVOLVIMENTO

Rodar comando 

```bash
docker compose up -d --build
```


**Acesse a aplicação:**

   Acesse a aplicação no endereço `http://localhost:8080`.

### Exemplos de Requisição

- **Para criar leilão:**

   ```sh
   curl -X POST http://localhost:8080/auction \
   -H "Content-Type: application/json" \
   -d '{
        "product_name": "Product",
        "category": "Product",
        "description": "A short description about",
        "condition": 1
       }'
   ```

- **Para buscar leilão por ID:**

   ```sh
   curl -X GET "http://localhost:8080/auction/{auctionId}" \
   -H "Content-Type: application/json"
   ```

- **Para buscar leilões por parâmetros de consulta:**

   ```sh
   curl -X GET "http://localhost:8080/auction?category=Product&condition=1&status=0" \
   -H "Content-Type: application/json"
   ```