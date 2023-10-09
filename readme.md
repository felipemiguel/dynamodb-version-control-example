
# DynamoDB Version Control Example

Este é um exemplo de projeto que demonstra como implementar controle de versão ao atualizar registros no DynamoDB usando Golang e LocalStack.

## Como Executar

Siga as etapas abaixo para executar o projeto:

1. Clone este repositório para o seu ambiente local:

   ```bash
   git clone https://github.com/felipemiguel/dynamodb-version-control-example.git
   cd dynamodb-version-control-example
   ```

2. Certifique-se de ter o Docker e o Docker Compose instalados.

3. Execute o script `startup.sh` para configurar e iniciar o ambiente LocalStack e criar a tabela do DynamoDB:

   ```bash
   ./startup.sh
   ```

## Artigo Relacionado
[Link para o Artigo no Medium](https://medium.com/@fee_miguel/implementando-controle-de-vers%C3%A3o-para-atualiza%C3%A7%C3%B5es-concorrentes-no-dynamodb-com-golang-3907744540b3)
