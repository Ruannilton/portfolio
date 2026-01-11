# scripts/seed_data.py
import argparse
import random
from datetime import datetime, timedelta, timezone
from typing import Any, Dict, List, Optional

import requests

SENIORITY_OPTIONS = ["JUNIOR", "MID_LEVEL", "SENIOR", "LEAD", "PRINCIPAL", "STAFF"]
LOCATION_OPTIONS = ["ON_SITE", "REMOTE", "HYBRID", "ANY"]
CONTRACT_OPTIONS = ["PJ", "CLT", "FREELANCER", "CONTRACTOR"]
CURRENCY_OPTIONS = ["BRL", "USD", "EUR", "GBP"]

# --- DIVERSIDADE EXPANDIDA ---
SKILL_POOL = [
    # Linguagens
    "Go", "Python", "TypeScript", "React", "Docker", "Kubernetes",
    "PostgreSQL", "Redis", "AWS", "Terraform", "GraphQL", "gRPC",
    "Node.js", "Vue.js", "Angular", "Java", "C#", "Ruby", "PHP",
    "MongoDB", "Elasticsearch", "RabbitMQ", "Svelte", "Next.js", "Django",
    "Rust", "Kotlin", "Swift", "C++", "Elixir", "Scala", "Dart", "Haskell",
    # Frameworks & Libs
    "FastAPI", "Spring Boot", "Laravel", "NestJS", "TailwindCSS", "Bootstrap",
    "Flutter", "React Native", ".NET Core", "Flask", "Express",
    # Infra & Tools
    "Ansible", "Jenkins", "GitHub Actions", "Prometheus", "Grafana", "Linux",
    "Nginx", "Apache Kafka", "Cassandra", "DynamoDB", "Azure", "GCP",
    "Git", "Jira", "Figma", "Datadog", "Splunk",
    # AI & Data
    "Pandas", "NumPy", "PyTorch", "TensorFlow", "OpenAI API", "Scikit-learn"
]

PROJECT_NAMES = [
    "Atlas", "Zephyr", "Aurora", "Nebula", "Pulse", "Vertex", "Quark", "Solstice",
    "Orion", "Helix", "Nimbus", "Photon", "Vortex", "Echo", "Blaze", "Nova",
    "Chimera", "Odyssey", "Titan", "Chronos", "Aigis", "Hyperion", "Zenith", "Omega",
    "Spectre", "Phantom", "Mirage", "Equinox", "Polaris", "Sirius"
]

BIO_VARIATIONS = [
    "Profissional com experiência em múltiplas stacks e entrega de produtos de alto impacto.",
    "Apaixonado por tecnologia, inovação e desafios de escala.",
    "Especialista em backend, mas com forte atuação em DevOps e Cloud.",
    "Foco em qualidade de código, automação e cultura ágil.",
    "Entusiasta de open source e comunidades técnicas.",
    "Experiência internacional em projetos globais e equipes distribuídas.",
    "Transformando ideias em soluções digitais robustas.",
    "Mentor de equipes e multiplicador de conhecimento técnico.",
    "Atuação em startups e grandes empresas, sempre buscando inovação.",
    "Desenvolvedor full-stack com paixão por performance e UX.",
    "Engenheiro focado em resiliência e sistemas de alta disponibilidade.",
    "Arquiteto de software com viés para soluções simples e eficientes.",
    "Explorador de novas tecnologias, focado atualmente em IA Generativa.",
    "Especialista em modernização de legados e migração para nuvem."
]

PROJECT_DESCRIPTIONS = [
    "Projeto focado em soluções inovadoras para o mercado financeiro (FinTech).",
    "Plataforma de automação de processos empresariais.",
    "Sistema de recomendação inteligente usando machine learning.",
    "Ferramenta de monitoramento em tempo real para aplicações web.",
    "API de alta performance para integração de sistemas legados.",
    "Aplicativo mobile para gestão de tarefas colaborativas.",
    "Dashboard analítico com visualização de dados interativa.",
    "Infraestrutura como código para ambientes escaláveis.",
    "Chatbot inteligente integrado a múltiplos canais de atendimento.",
    "Sistema de autenticação e autorização robusto para microserviços.",
    "Plataforma de E-commerce com alta capacidade de tráfego.",
    "Solução IoT para monitoramento de sensores industriais.",
    "Marketplace de serviços integrando pagamentos digitais.",
    "Sistema de telemedicina compatível com normas de segurança de dados.",
    "Rede social corporativa para engajamento de times remotos."
]

EXPERIENCE_DESCRIPTIONS = [
    "Responsável pelo desenvolvimento de serviços escaláveis e APIs REST.",
    "Liderança técnica em squads multidisciplinares.",
    "Implementação de pipelines CI/CD e automação de deploy.",
    "Migração de sistemas monolíticos para arquitetura de microsserviços.",
    "Otimização de queries e modelagem de banco de dados relacional e NoSQL.",
    "Mentoria de desenvolvedores juniores e code review.",
    "Integração de sistemas com provedores de nuvem (AWS, GCP, Azure).",
    "Desenvolvimento de testes automatizados e prática de TDD.",
    "Participação em decisões de arquitetura e escolha de tecnologias.",
    "Atuação em projetos ágeis com entregas contínuas.",
    "Redução de custos de infraestrutura através de otimização de recursos.",
    "Implementação de observabilidade e logging centralizado.",
    "Desenvolvimento de interfaces responsivas e acessíveis."
]

ROLES = [
    "Backend Engineer", "Full-Stack Developer", "Tech Lead", "DevOps Engineer",
    "Frontend Developer", "Cloud Architect", "QA Engineer", "Product Owner",
    "Site Reliability Engineer (SRE)", "Data Engineer", "Mobile Developer",
    "Security Engineer", "Solutions Architect", "Staff Engineer", "Engineering Manager"
]
# --------------------------------

def random_subset(pool: List[str], min_items: int, max_items: int) -> List[str]:
    count = random.randint(min_items, max_items)
    return random.sample(pool, k=min(count, len(pool)))


def format_datetime(dt: datetime) -> str:
    """Formata datetime para RFC3339 (compatível com Go time.Time)."""
    return dt.strftime("%Y-%m-%dT%H:%M:%SZ")


def random_date_range(max_years_back: int = 8) -> tuple[str, Optional[str]]:
    start = datetime.now(timezone.utc) - timedelta(days=random.randint(180, 365 * max_years_back))
    if random.random() < 0.30: # Aumentei chance de ser emprego atual (sem data fim)
        end_iso = None
    else:
        end = start + timedelta(days=random.randint(180, 1000))
        # Garante que data fim não seja no futuro distante
        if end > datetime.now(timezone.utc):
            end = datetime.now(timezone.utc)
        end_iso = format_datetime(end)
    return format_datetime(start), end_iso


def build_experiences() -> List[Dict[str, Any]]:
    experiences = []
    # Mais variação na quantidade de experiências
    for _ in range(random.randint(1, 5)):
        start_iso, end_iso = random_date_range()
        experiences.append({
            "company": f"{random.choice(['Acme', 'Globex', 'Initech', 'Stark', 'Umbrella', 'Wayne', 'Cyberdyne', 'Massive', 'Hooli', 'Pied Piper'])} {random.choice(['Labs', 'Inc', 'Corp', 'Systems', 'Tech'])}",
            "role": random.choice(ROLES),
            "startDate": start_iso,
            "endDate": end_iso,
            "description": random.choice(EXPERIENCE_DESCRIPTIONS),
            "techStack": random_subset(SKILL_POOL, 3, 8),
        })
    return experiences


def build_projects() -> List[Dict[str, Any]]:
    projects = []
    for _ in range(random.randint(1, 4)):
        name = f"Project {random.choice(PROJECT_NAMES)}"
        projects.append({
            "name": name,
            "description": random.choice(PROJECT_DESCRIPTIONS),
            "repoUrl": f"https://github.com/seed/{name.replace(' ', '-').lower()}",
            "liveUrl": f"https://{name.replace(' ', '').lower()}.example.com",
            "tags": random_subset(SKILL_POOL, 2, 6),
        })
    return projects


def build_educations() -> List[Dict[str, Any]]:
    educations = []
    for _ in range(random.randint(1, 2)): # Pelo menos 1 educação
        start_iso, end_iso = random_date_range(max_years_back=12)
        # Adicionei mais universidades, incluindo contexto BR
        inst = random.choice([
            "USP", "UNICAMP", "UFPE", "PUC-Rio", "UFRJ", "UFMG", "UTFPR",
            "ITA", "UFRGS", "UNB", "FIAP", "Harvard", "MIT", "Stanford"
        ])
        educations.append({
            "institution": inst,
            "degree": random.choice(["Bacharelado", "Mestrado", "Doutorado", "Tecnólogo", "MBA"]),
            "field": random.choice([
                "Ciência da Computação", "Engenharia de Software", "Sistemas de Informação", 
                "Engenharia Elétrica", "Análise de Sistemas", "Matemática Computacional"
            ]),
            "startDate": start_iso,
            "endDate": end_iso,
        })
    return educations


def build_profile_payload(seed_index: int) -> Dict[str, Any]:
    headline = random.choice([
        "Desenvolvedor apaixonado por performance",
        "Especialista em plataformas distribuídas",
        "Engenheiro focado em produtos escaláveis",
        "Tech Lead com experiência em equipes globais",
        "DevOps focado em automação e cloud",
        "Arquiteto de soluções inovadoras",
        "Frontend expert em UX e acessibilidade",
        "Backend lover com paixão por APIs rápidas",
        "Full-stack Developer | Golang | React",
        "Cloud Native Engineer & Kubernetes Enthusiast"
    ])
    
    return {
        "headline": headline,
        "bio": random.choice(BIO_VARIATIONS),
        "seniority": random.choice(SENIORITY_OPTIONS),
        "years_of_experience": random.randint(1, 20),
        "open_to_work": random.random() < 0.7,
        "salary_expectation": round(random.uniform(4000, 35000), 2),
        "currency": random.choice(CURRENCY_OPTIONS),
        "contract_type": random.choice(CONTRACT_OPTIONS),
        "location": random.choice(LOCATION_OPTIONS),
        "remote_only": random.random() < 0.5,
        "skills": random_subset(SKILL_POOL, 4, 12),
        "social_links": {
            "linkedin": f"https://www.linkedin.com/in/userseed{seed_index}",
            "github": f"https://github.com/userseed{seed_index}",
            "website": f"https://userseed{seed_index}.dev",
        },
        "experiences": build_experiences(),
        "projects": build_projects(),
        "educations": build_educations(),
    }


def register_user(session: requests.Session, base_url: str, email: str, password: str, index: int) -> None:
    payload = {
        "firstName": "Seed",
        "lastName": f"User{index}",
        "email": email,
        "password": password,
    }
    try:
        response = session.post(f"{base_url}/auth/register", json=payload, timeout=15)
        # 409 Conflict é aceitável se o usuário já existir, seguimos para o login
        if response.status_code not in (201, 409):
            print(f"⚠️  Register failed for {email} ({response.status_code}): {response.text}")
    except Exception as e:
        print(f"❌ Error registering {email}: {e}")


def login_user(session: requests.Session, base_url: str, email: str, password: str) -> str:
    payload = {"email": email, "password": password}
    response = session.post(f"{base_url}/auth/login", json=payload, timeout=15)
    if response.status_code != 200:
        raise RuntimeError(f"Login failed ({response.status_code}): {response.text}")
    data = response.json()
    if "access_token" not in data:
        raise RuntimeError("Login response missing access_token")
    return data["access_token"]


def create_profile(session: requests.Session, base_url: str, token: str, profile_payload: Dict[str, Any]) -> None:
    headers = {"Authorization": f"Bearer {token}"}
    response = session.post(f"{base_url}/portfolio/", json=profile_payload, headers=headers, timeout=15)
    
    if response.status_code in (201, 409):
        return
    if "already exists" in response.text.lower():
        return
        
    print(f"⚠️  Profile creation failed ({response.status_code}): {response.text}")


def seed(base_url: str, count: int, start_index: int) -> None:
    session = requests.Session()
    password = "userseed123"

    print(f"🚀 Starting seed process...")
    print(f"Target: {base_url}")
    print(f"Creating {count} users starting from index {start_index} (userseed{start_index}@email.com)")

    success_count = 0
    
    # Loop ajustado para usar start_index
    for idx in range(start_index, start_index + count):
        email = f"userseed{idx}@email.com"
        try:
            print(f"Processing {email}...", end=" ", flush=True)
            register_user(session, base_url, email, password, idx)
            token = login_user(session, base_url, email, password)
            profile_payload = build_profile_payload(idx)
            create_profile(session, base_url, token, profile_payload)
            print("✅ OK")
            success_count += 1
        except Exception as e:
            print(f"❌ Failed: {str(e)}")
    
    print(f"\n✨ Seeding complete! Successfully processed {success_count}/{count} users.")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Seed users and portfolios through the API.")
    parser.add_argument("--base-url", default="http://localhost:8080", help="API base URL (default: http://localhost:8080)")
    
    # Flags solicitadas
    parser.add_argument("--count", type=int, default=10, help="Number of users to create")
    parser.add_argument("--start-index", type=int, default=1, help="Starting index number for user emails (useful to avoid collisions)")
    
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    seed(args.base_url.rstrip("/"), args.count, args.start_index)


if __name__ == "__main__":
    main()