/**
 * portfolio-editor.js - Funções do editor de portfólio
 */

function showCreateForm() {
    document.getElementById('portfolio-edit').classList.remove('hidden');
    document.getElementById('portfolio-empty').classList.add('hidden');
}

function toggleEditMode() {
    const viewMode = document.getElementById('portfolio-view');
    const editMode = document.getElementById('portfolio-edit');
    const emptyMode = document.getElementById('portfolio-empty');
    if (viewMode) viewMode.classList.toggle('hidden');
    if (editMode) editMode.classList.toggle('hidden');
    if (!viewMode && emptyMode && editMode.classList.contains('hidden')) {
        emptyMode.classList.remove('hidden');
    }
}

function removeItem(button, type) {
    button.closest(`.${type}-item`).remove();
}

function addExperience() {
    const container = document.getElementById('experiences-container');
    const html = `
    <div class="experience-item border border-gray-200 rounded-lg p-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Empresa</label>
                <input type="text" data-field="company" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Cargo</label>
                <input type="text" data-field="role" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Data Início</label>
                <input type="date" data-field="startDate" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Data Fim</label>
                <input type="date" data-field="endDate" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Descrição</label>
                <textarea data-field="description" rows="2" class="w-full p-2 border border-gray-300 rounded-lg"></textarea>
            </div>
            <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Tech Stack (separado por vírgula)</label>
                <input type="text" data-field="techStack" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
        </div>
        <button type="button" onclick="removeItem(this, 'experience')" class="mt-2 text-red-600 hover:text-red-800 text-sm">🗑️ Remover</button>
    </div>
`;
    container.insertAdjacentHTML('beforeend', html);
}

function escapeHtml(str) {
    if (!str) return '';
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

function updateBio(content) {
    const textArea = document.getElementById('bio-textarea');
    textArea.value = escapeHtml(content);
}

function updateSocialLinks(links) {
    document.getElementById('social_links.linkedin').value = links.linkedin || '';
    document.getElementById('social_links.github').value = links.github || '';
    document.getElementById('social_links.website').value = links.website || '';
}

function addProject(data = null) {
    const container = document.getElementById('projects-container');

    // Define valores padrão
    const name = data ? data.name : '';
    const description = data ? data.description : '';
    const repoUrl = data ? data.repoUrl : '';
    const liveUrl = data ? data.liveUrl : '';
    const tags = data ? (data.tags || []).join(', ') : '';

    // NOVOS VALORES
    const provider = data ? (data.provider || '') : '';
    const providerId = data ? (data.providerId || '') : '';

    const html = `
    <div class="project-item border border-gray-200 rounded-lg p-4 animate-fade-in">
        
        <input type="hidden" data-field="provider" value="${escapeHtml(provider)}">
        <input type="hidden" data-field="providerId" value="${escapeHtml(providerId)}">

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nome</label>
                <input type="text" data-field="name" value="${escapeHtml(name)}" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Tags (separadas por vírgula)</label>
                <input type="text" data-field="tags" value="${escapeHtml(tags)}" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div class="md:col-span-2">
                <label class="block text-sm font-medium text-gray-700 mb-1">Descrição</label>
                <textarea data-field="description" rows="2" class="w-full p-2 border border-gray-300 rounded-lg">${escapeHtml(description)}</textarea>
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">URL do Repositório</label>
                <input type="url" data-field="repoUrl" value="${repoUrl}" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">URL Demo</label>
                <input type="url" data-field="liveUrl" value="${liveUrl}" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
        </div>
        <button type="button" onclick="removeItem(this, 'project')" class="mt-2 text-red-600 hover:text-red-800 text-sm">🗑️ Remover</button>
    </div>
`;
    container.insertAdjacentHTML('beforeend', html);
}

function addEducation() {
    const container = document.getElementById('educations-container');
    const html = `
    <div class="education-item border border-gray-200 rounded-lg p-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Instituição</label>
                <input type="text" data-field="institution" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Grau</label>
                <input type="text" data-field="degree" class="w-full p-2 border border-gray-300 rounded-lg" placeholder="Bacharelado, Mestrado...">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Área</label>
                <input type="text" data-field="field" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Data Início</label>
                <input type="date" data-field="startDate" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
            <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Data Fim</label>
                <input type="date" data-field="endDate" class="w-full p-2 border border-gray-300 rounded-lg">
            </div>
        </div>
        <button type="button" onclick="removeItem(this, 'education')" class="mt-2 text-red-600 hover:text-red-800 text-sm">🗑️ Remover</button>
    </div>
`;
    container.insertAdjacentHTML('beforeend', html);
}

function prepareFormData(event) {
    event.preventDefault();
    const form = event.target.closest('form');

    const data = {
        headline: form.querySelector('[name="headline"]').value,
        bio: form.querySelector('[name="bio"]').value,
        seniority: form.querySelector('[name="seniority"]').value,
        years_of_experience: parseInt(form.querySelector('[name="years_of_experience"]').value) || 0,
        open_to_work: form.querySelector('[name="open_to_work"]').checked,
        salary_expectation: parseFloat(form.querySelector('[name="salary_expectation"]').value) || 0,
        currency: form.querySelector('[name="currency"]').value,
        contract_type: form.querySelector('[name="contract_type"]').value,
        location: form.querySelector('[name="location"]').value,
        remote_only: form.querySelector('[name="remote_only"]').checked,
        skills: document.getElementById('skills-input').value.split(',').map(s => s.trim()).filter(s => s),
        social_links: {
            linkedin: form.querySelector('[name="social_links.linkedin"]').value,
            github: form.querySelector('[name="social_links.github"]').value,
            website: form.querySelector('[name="social_links.website"]').value
        },
        experiences: [],
        projects: [],
        educations: []
    };

    // Collect experiences
    document.querySelectorAll('.experience-item').forEach(item => {
        const exp = {
            company: item.querySelector('[data-field="company"]').value,
            role: item.querySelector('[data-field="role"]').value,
            startDate: item.querySelector('[data-field="startDate"]').value ? new Date(item.querySelector('[data-field="startDate"]').value).toISOString() : null,
            endDate: item.querySelector('[data-field="endDate"]').value ? new Date(item.querySelector('[data-field="endDate"]').value).toISOString() : null,
            description: item.querySelector('[data-field="description"]').value,
            techStack: item.querySelector('[data-field="techStack"]').value.split(',').map(s => s.trim()).filter(s => s)
        };
        if (exp.company || exp.role) data.experiences.push(exp);
    });

    // Collect projects
    document.querySelectorAll('.project-item').forEach(item => {
        const providerVal = item.querySelector('[data-field="provider"]').value;
        const providerIdVal = item.querySelector('[data-field="providerId"]').value;

        const proj = {
            name: item.querySelector('[data-field="name"]').value,
            description: item.querySelector('[data-field="description"]').value,
            repoUrl: item.querySelector('[data-field="repoUrl"]').value,
            liveUrl: item.querySelector('[data-field="liveUrl"]').value,
            tags: item.querySelector('[data-field="tags"]').value.split(',').map(s => s.trim()).filter(s => s),
            provided: providerVal || null,
            providerId: providerIdVal || null
        };
        if (proj.name) data.projects.push(proj);
    });

    // Collect educations
    document.querySelectorAll('.education-item').forEach(item => {
        const edu = {
            institution: item.querySelector('[data-field="institution"]').value,
            degree: item.querySelector('[data-field="degree"]').value,
            field: item.querySelector('[data-field="field"]').value,
            startDate: item.querySelector('[data-field="startDate"]').value ? new Date(item.querySelector('[data-field="startDate"]').value).toISOString() : null,
            endDate: item.querySelector('[data-field="endDate"]').value ? new Date(item.querySelector('[data-field="endDate"]').value).toISOString() : null
        };
        if (edu.institution) data.educations.push(edu);
    });

    // Send via fetch with JSON
    fetch('/portfolio/html', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
        .then(response => response.text())
        .then(html => {
            document.getElementById('portfolio-view').innerHTML = html;
            // Re-execute scripts in the new HTML
            const scripts = document.getElementById('portfolio-container').querySelectorAll('script');
            scripts.forEach(script => {
                const newScript = document.createElement('script');
                newScript.textContent = script.textContent;
                document.body.appendChild(newScript);
                document.body.removeChild(newScript);
            });
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Erro ao atualizar portfólio');
        });
}

/**
 * Importa do GitHub verificando Provider e ProviderId
 */
async function importGithubProjects() {
    const btn = document.getElementById('btn-import-github');
    const originalText = btn.innerHTML;

    btn.disabled = true;
    btn.innerHTML = `Importando...`;

    try {
        // 1. Mapear projetos existentes usando uma chave única "PROVIDER:ID"
        const existingKeys = new Set();

        document.querySelectorAll('.project-item').forEach(item => {
            const p = item.querySelector('input[data-field="provider"]').value;
            const pid = item.querySelector('input[data-field="providerId"]').value;

            // Só adiciona ao Set se tiver ID (projetos manuais não têm ID externo)
            if (p && pid) {
                existingKeys.add(`${p}:${pid}`);
            }
        });

        // 2. Buscar da API
        const response = await fetch('/sync/github'); // ou /api/github dependendo da sua rota
        if (!response.ok) throw new Error('Falha na API');
        const data = await response.json();

        let addedCount = 0;

        if (data.repositories && Array.isArray(data.repositories)) {
            data.repositories.forEach(repo => {
                // Monta a chave do projeto que veio da API
                // Assumindo Provider "GITHUB" fixo para essa importação
                const currentKey = `GITHUB:${repo.ProviderId}`;

                // VERIFICAÇÃO DE DUPLICIDADE
                if (!existingKeys.has(currentKey)) {
                    addProject({
                        name: repo.Name,
                        description: repo.Description,
                        repoUrl: repo.Url,
                        liveUrl: repo.Homepage || '', // Garanta que seu backend retorna Homepage se quiser usar
                        tags: [...(repo.Languages || []), ...(repo.Topics || [])],
                        provider: 'GITHUB',       // <--- Fixo
                        providerId: repo.ProviderId // <--- Vindo do backend
                    });
                    addedCount++;
                }
            });
        }

        updateBio(data.bio || '');
        const socialLinks = {
            linkedin: data.linkedinUrl ? data.linkedinUrl : '',
            github: data.githubUrl ? data.githubUrl : '',
            website: data.genericUrl ? data.genericUrl : ''
        };
        updateSocialLinks(socialLinks);
        if (addedCount > 0) {
            alert(`${addedCount} projetos novos adicionados!`);
        } else {
            alert('Todos os projetos do GitHub já estão no seu portfólio.');
        }

    } catch (error) {
        console.error(error);
        alert('Erro: ' + error.message);
    } finally {
        btn.disabled = false;
        btn.innerHTML = originalText;
    }
}