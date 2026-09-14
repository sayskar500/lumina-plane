# 🤖 Lumina-Plane Agentic Flow Documentation

This document serves as the living blueprint for the agentic behaviors, orchestration patterns, and skill mappings within the Lumina-Plane platform. Following enterprise-grade AI engineering standards, we treat agent behavior as an architectural component rather than just a prompt.

## 🏛️ Architectural Pattern: Augmented LLM (Single Agent)
Currently, Lumina-Plane operates as an **Augmented LLM** system. It uses a centralized AI Gateway that augments raw user input with project-specific context (via the Prompt Registry) before routing to a high-performance LLM.

### 🔄 Logic Flow
`User Request` $\rightarrow$ `CLI` $\rightarrow$ `AI Gateway (Server)` $\rightarrow$ `Prompt Registry (MongoDB)` $\rightarrow$ `Context Augmentation` $\rightarrow$ `LLM (Groq/Gemini)` $\rightarrow$ `Response` $\rightarrow$ `Token Logging (MongoDB)`

### 🛠️ Skill Map (Capabilities)
| Skill | Input | Logic | Output |
| :--- | :--- | :--- | :--- |
| **Health Check** | - | Verifies API $\rightarrow$ DB connectivity | Connectivity Status |
| **Context Augmentation** | `ProjectID`, `UserPrompt` | Fetches `IsActive=true` template for `ProjectID` and prepends to prompt | Augmented Prompt |
| **AI Generation** | `AugmentedPrompt`, `Model` | Routes to LLM provider via API Gateway | AI Response + Token Usage |
| **Usage Tracking** | `ProjectID`, `Tokens`, `Model` | Asynchronous write to `token_logs` collection | Success/Failure |
| **Prompt Management** | `ProjectID`, `Template` | Deactivates existing active prompts $\rightarrow$ Inserts new version | New Version ID |

## 📈 Future Agentic Evolutions
As the platform scales, we will migrate from a single augmented LLM to more complex patterns:

### 1. Sequential Pipeline (Compliance & Audit)
For infrastructure provisioning, we will implement a *Plan $\rightarrow$ Review $\rightarrow$ Execute* pipeline.
- **Planner Agent**: Decomposes the provisioning request into Terraform steps.
- **Reviewer Agent**: Validates the plan against security policies.
- **Executor Agent**: Applies the Terraform plan.

### 2. Evaluator-Optimizer (Quality Assurance)
For the Prompt Registry, we will implement an iterative refinement loop.
- **Generator**: Proposes a new prompt template.
- **Evaluator**: Runs a set of test cases and scores the output.
- **Optimizer**: Refines the template based on the score.

### 3. Hierarchical Orchestrator (Multi-Agent)
For complex developer requests (e.g., "Setup a new project with MongoDB and a Go server"), a supervisor agent will delegate to:
- **Infra Agent**: Handles Terraform/Cloud resources.
- **Code Agent**: Generates boilerplate scaffolding.
- **Ops Agent**: Configures K8s manifests.

## 📏 Documentation Standard for Modifications
Every modification to the system's behavior must be updated here using the following template:

### [Feature Name]
- **Pattern Change**: (e.g., "Single Agent" $\rightarrow$ "Sequential Pipeline")
- **New Skills Added**: List of new tools/functions provided to the agent.
- **Logic Update**: Description of how the reasoning loop has changed.
- **Memory/State Change**: Any new data being persisted to MongoDB.
- **Flow Diagram**: Mermaid diagram updated to reflect the new path.
