# 🍽️ Sistema de Restaurante Full-Stack (Go + SQLite)

Um sistema simples e funcional para gestão de pedidos em restaurantes. O projeto conta com uma interface para o cliente realizar pedidos e um painel para a cozinha gerenciar a produção em tempo real.



## 🚀 Funcionalidades

- **Cardápio Digital:** Listagem de pratos com descrição e preço.
- **Realização de Pedidos:** O cliente informa nome, mesa e seleciona os pratos.
- **Painel da Cozinha:** Visualização de pedidos pendentes em tempo real (atualização automática).
- **Gestão de Status:** Possibilidade de marcar pedidos como concluídos.
- **Histórico e Limpeza:** Área para visualizar pedidos finalizados e opção para limpar o histórico do banco de dados.

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** [Go (Golang)](https://go.dev/)
- **Framework Web:** [Gin Gonic](https://gin-gonic.com/)
- **Banco de Dados:** SQLite (via [GORM](https://gorm.io/))
- **Frontend:** HTML5, CSS3 e JavaScript (Vanilla)

## 📂 Estrutura do Projeto

```text
├── main.go          # Servidor backend e API
├── index.html       # Interface do Cliente (Cardápio)
├── cozinha.html     # Interface do Restaurante (Cozinha)
├── restaurant.db    # Banco de dados SQLite (gerado automaticamente)
└── go.mod           # Gerenciador de dependências Go
```
## ⚙️ Como Executar Localmente
  Pré-requisito
  
  Ter o Go instalado (versão 1.20 ou superior recomendada).
  
## Passo a Passo
- Clone o repositório:
  git clone https://github.com/Cley-gabriel/Meu-Restaurante
  cd Meu-Restaurante

- Instale as dependências:
  go mod tidy
  
  # Baixa o framework da API
  go get github.com/gin-gonic/gin

  # Baixa o suporte para conexões de outros sites (CORS)
  go get github.com/gin-contrib/cors

  # Baixa o ORM (trabalha com banco de dados)
  go get gorm.io/gorm

  # Baixa o driver específico do SQLite que não precisa de C++ 
  go get github.com/glebarez/sqlite

- Execute o servidor:
  
  go run main.go
  
- Acesse no seu navegador:
 
  Cliente: http://localhost:8080
  
  Cozinha: http://localhost:8080/cozinha

  <img width="678" height="607" alt="image" src="https://github.com/user-attachments/assets/c5aa22cf-7679-48ed-881f-bb96646d8d33" />
  <img width="1282" height="601" alt="image" src="https://github.com/user-attachments/assets/f5584b1e-f674-4e87-9d87-7969f00c7a9f" />

  


  
