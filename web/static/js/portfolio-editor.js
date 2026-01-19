/**
 * portfolio-editor.js - Funções do editor de portfólio
 */

pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/2.16.105/pdf.worker.min.js';

// htmx.defineExtension('submit-json', {
//     onEvent: function (name, evt) {
//         if (name === "htmx:configRequest") {
//             evt.detail.headers['Content-Type'] = "application/json";
//         }
//     },
//     encodeParameters: function(xhr, parameters, elt) {
//         xhr.overrideMimeType('text/json');
//         // Garante que o JSON seja gerado exatamente como o objeto está,
//         // sem tentar converter valores para string antes.
//         return JSON.stringify(parameters);
//     }
// });

// document.body.addEventListener('htmx:configRequest', function(evt) {
//     // Verifica se é o formulário de perfil pelo endpoint ou ID
//     if (evt.target.getAttribute('hx-put') === '/app/profile') {
        
//         // 1. Pega o elemento formulário
//         const form = evt.target;
        
//         // 2. Chama sua função para gerar o JSON limpo
//         const complexData = prepareFormData(form);
        
//         // 3. Sobrescreve os parâmetros que o HTMX enviaria.
//         // Como você está usando hx-ext="json-enc", ele vai pegar esse objeto
//         // e serializar automaticamente para JSON no corpo da requisição.
//         evt.detail.parameters = complexData;
        
//         // Debug opcional: verifique no console o que está indo
//         console.log("Enviando JSON customizado:", complexData);
//     }
// });

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


function sendFormData(data){
     // Send via fetch with JSON
    fetch('/app/profile', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    })
        .then(response => {
            if (!response.ok) throw new Error('Falha na requisição');
            return response.text();
        })
        .then(html => {
            // 1. Atualiza o conteúdo da View
            const viewElement = document.getElementById('portfolio-view');
            viewElement.innerHTML = html;

            // 2. Re-executa scripts que vieram no HTML novo (buscando dentro do viewElement)
            const scripts = viewElement.querySelectorAll('script');
            scripts.forEach(script => {
                const newScript = document.createElement('script');
                newScript.textContent = script.textContent;
                document.body.appendChild(newScript);
                document.body.removeChild(newScript);
            });

            // 3. Volta para o modo de visualização (esconde o form, mostra o view)
            // Se você já estiver no modo edit, chamar toggleEditMode deve inverter.
            // Certifique-se que esta função faz o que você espera, ou force as classes aqui.
            toggleEditMode(); 
        })
        .catch(error => {
            console.error(error);
            alert('Erro ao atualizar portfólio');
        });
}



function prepareFormData(form) {
   console.log(form);
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

    // Coletar Experiências
    form.querySelectorAll('.experience-item').forEach(item => {
        const exp = {
            company: item.querySelector('[data-field="company"]').value,
            role: item.querySelector('[data-field="role"]').value,
            // Cuidado com timezone no toISOString. Se for apenas data, considere enviar a string crua "YYYY-MM-DD"
            startDate: item.querySelector('[data-field="startDate"]').value ? new Date(item.querySelector('[data-field="startDate"]').value).toISOString() : null,
            endDate: item.querySelector('[data-field="endDate"]').value ? new Date(item.querySelector('[data-field="endDate"]').value).toISOString() : null,
            description: item.querySelector('[data-field="description"]').value,
            techStack: item.querySelector('[data-field="techStack"]').value.split(',').map(s => s.trim()).filter(s => s)
        };
        // Pequena validação para não enviar objetos vazios
        if (exp.company || exp.role) data.experiences.push(exp);
    });

    // Coletar Projetos
    form.querySelectorAll('.project-item').forEach(item => {
        const providerVal = item.querySelector('[data-field="provider"]').value;
        const providerIdVal = item.querySelector('[data-field="providerId"]').value;

        const proj = {
            name: item.querySelector('[data-field="name"]').value,
            description: item.querySelector('[data-field="description"]').value,
            // ATENÇÃO AQUI: Verifique se sua Struct Go espera "repoUrl" ou "repo_url"
            repoUrl: item.querySelector('[data-field="repoUrl"]').value,
            liveUrl: item.querySelector('[data-field="liveUrl"]').value,
            tags: item.querySelector('[data-field="tags"]').value.split(',').map(s => s.trim()).filter(s => s),
            provider: providerVal || null,
            providerId: providerIdVal || null
        };
        if (proj.name) data.projects.push(proj);
    });

    // Coletar Educação
    form.querySelectorAll('.education-item').forEach(item => {
        const edu = {
            institution: item.querySelector('[data-field="institution"]').value,
            degree: item.querySelector('[data-field="degree"]').value,
            field: item.querySelector('[data-field="field"]').value,
            startDate: item.querySelector('[data-field="startDate"]').value ? new Date(item.querySelector('[data-field="startDate"]').value).toISOString() : null,
            endDate: item.querySelector('[data-field="endDate"]').value ? new Date(item.querySelector('[data-field="endDate"]').value).toISOString() : null
        };
        if (edu.institution) data.educations.push(edu);
    });
    console.log(data)
   return data;
}

function submitProfileForm(event) {
    event.preventDefault();
    const form = event.target;
    const data = prepareFormData(form);
    sendFormData(data);
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

async function handleLinkedinUploadEvent(event) {
    const file = event.target.files[0];
    if (!file) return;
    try {
        const arrayBuffer = await file.arrayBuffer();
        const pdf = await pdfjsLib.getDocument(arrayBuffer).promise;
        let fullText = "";

        for (let i = 1; i <= pdf.numPages; i++) {
            const page = await pdf.getPage(i);
            const textContent = await page.getTextContent();
            fullText += textContent.items.map(item => item.str).join('\n') + "\n";
        }

        const profileData = parseLinkedInProfile(fullText);
        importLinkedInData(profileData);
        console.log(profileData);

    } catch (error) {
        console.error(error);
    }
}

function appendSkills(skills) {
    const skillsInput = document.getElementById('skills-input');
    const existingSkills = skillsInput.value.split(',').map(s => s.trim().toLowerCase());
    skills.forEach(skill => {
        if (!existingSkills.includes(skill.toLowerCase())) {
            existingSkills.push(skill);
        }
    });
    skillsInput.value = existingSkills.join(', ');
}


function importLinkedInData(profile) {
    updateBio(profile.summary || '');
    const socialLinks = {
        linkedin: profile.contact.linkedin || '',
        github: profile.contact.github || '',
        website: ''
    };
    updateSocialLinks(socialLinks);
    appendSkills(profile.skills || []);
    profile.experience.forEach(exp => {
        addExperience();
        const expItems = document.querySelectorAll('.experience-item');
        const lastExp = expItems[expItems.length - 1];
        lastExp.querySelector('[data-field="company"]').value = exp.company || '';
        lastExp.querySelector('[data-field="role"]').value = exp.role || '';
        if (exp.start) lastExp.querySelector('[data-field="startDate"]').value = exp.start.toISOString().split('T')[0];
        if (exp.end) lastExp.querySelector('[data-field="endDate"]').value = exp.end.toISOString().split('T')[0];
        lastExp.querySelector('[data-field="description"]').value = exp.description || '';
    });
    profile.education.forEach(edu => {
        addEducation();
        const eduItems = document.querySelectorAll('.education-item');
        const lastEdu = eduItems[eduItems.length - 1];
        lastEdu.querySelector('[data-field="institution"]').value = edu.school || '';
        lastEdu.querySelector('[data-field="degree"]').value = edu.degree || '';
        if (edu.start) lastEdu.querySelector('[data-field="startDate"]').value = edu.start.toISOString().split('T')[0];
        if (edu.end) lastEdu.querySelector('[data-field="endDate"]').value = edu.end.toISOString().split('T')[0];
    });

}

const monthMap = {
    'janeiro': 0, 'fevereiro': 1, 'março': 2, 'abril': 3, 'maio': 4, 'junho': 5,
    'julho': 6, 'agosto': 7, 'setembro': 8, 'outubro': 9, 'novembro': 10, 'dezembro': 11,
    'jan': 0, 'feb': 1, 'mar': 2, 'apr': 3, 'may': 4, 'jun': 5,
    'jul': 6, 'aug': 7, 'sep': 8, 'oct': 9, 'nov': 10, 'dec': 11
};

function parsePTDate(dateStr) {
    if (!dateStr) return null;
    const str = dateStr.trim().toLowerCase();

    const match = str.match(/([a-zç]+)\.?\s+(?:de\s+)?(\d{4})/);
    if (match) {
        const monthIndex = monthMap[match[1]];
        const year = parseInt(match[2], 10);
        if (monthIndex !== undefined && !isNaN(year)) return new Date(year, monthIndex, 1);
    }
    return null;
}

function parseDateRange(rawDateString) {
    const parts = rawDateString.split(/\s+-\s+/);
    return {
        start: parsePTDate(parts[0]),
        end: parts[1] ? parsePTDate(parts[1]) : null
    };
}

function parseLinkedInProfile(rawText) {
    const lines = rawText
        .split('\n')
        .map(line => line.trim())
        .filter(line => line.length > 0)
        .filter(line => !/^Page\s+\d+\s+of\s+\d+$/i.test(line) && line !== "Page");

    const profile = {
        contact: { email: "", linkedin: "", github: "" },
        name: "",
        headline: "",
        summary: "",
        skills: [],
        languages: [],
        experience: [],
        education: []
    };

    const findIndex = (keyword) => lines.findIndex(l => l.toLowerCase() === keyword.toLowerCase());

    const extractSection = (startKey, endKeys) => {
        const start = findIndex(startKey);
        if (start === -1) return [];
        let end = lines.length;
        for (const key of endKeys) {
            const idx = lines.findIndex((l, i) => i > start && l.toLowerCase().startsWith(key.toLowerCase()));
            if (idx !== -1 && idx < end) end = idx;
        }
        return lines.slice(start + 1, end);
    };

    // ==========================================================
    // CORREÇÃO APLICADA AQUI (Extração de Contato)
    // ==========================================================

    // 1. Email
    const emailMatch = rawText.match(/[\w.-]+@[\w.-]+\.\w+/);
    if (emailMatch) profile.contact.email = emailMatch[0];

    // 2. LinkedIn e Github (Iterando linhas para pegar quebras)
    for (let i = 0; i < lines.length; i++) {
        const line = lines[i];

        // LINKEDIN
        if (line.includes("linkedin.com")) {
            let fullUrl = line.trim();

            // Verifica se a próxima linha parece ser continuação
            // Ex: termina com "-" OU a próxima linha tem "(LinkedIn)"
            if (lines[i + 1] && (fullUrl.endsWith("-") || lines[i + 1].includes("(LinkedIn)"))) {
                // Pega a próxima linha e remove o rótulo "(LinkedIn)"
                const continuation = lines[i + 1].replace("(LinkedIn)", "").trim();
                fullUrl += continuation;
                // Opcional: pular a próxima linha no loop se necessário, mas não crítico aqui
            }
            profile.contact.linkedin = fullUrl;
        }

        // GITHUB
        if (line.includes("github.com")) {
            // Mesma lógica se o github estiver quebrado e a proxima linha for "(Portfolio)"
            let fullUrl = line.trim();
            if (lines[i + 1] && lines[i + 1].includes("(Portfolio)")) {
                const continuation = lines[i + 1].replace("(Portfolio)", "").trim();
                fullUrl += continuation;
            }
            profile.contact.github = fullUrl;
        }
    }

    // --- Resto da lógica (Nome, Resumo, etc...) ---

    const resumoIndex = findIndex("Resumo");
    if (resumoIndex > 2) {
        let foundHeadline = false;
        for (let i = resumoIndex - 1; i >= 0; i--) {
            const line = lines[i];
            if (line.includes("Brasil") || line.includes("Bahia") || line.includes("Janeiro") || line.includes("Paulo")) continue;
            if (!foundHeadline) {
                profile.headline = line;
                foundHeadline = true;
            } else if (!line.includes("@") && !line.includes("www") && !line.includes("Contato")) {
                profile.name = line;
                break;
            }
        }
    }

    profile.summary = extractSection("Resumo", ["Experiência", "Principais competências"]).join(" ");
    profile.skills = extractSection("Principais competências", ["Languages", "Certifications", "Resumo", "Experiência"]);

    const langLines = extractSection("Languages", ["Certifications", "Publications", "Resumo"]);
    for (let i = 0; i < langLines.length; i++) {
        if (langLines[i + 1] && langLines[i + 1].includes("(")) {
            profile.languages.push(`${langLines[i]} ${langLines[i + 1]}`);
            i++;
        } else profile.languages.push(langLines[i]);
    }

    // Experiência
    const expLines = extractSection("Experiência", ["Formação acadêmica", "Education"]);
    const dateRegex = /([a-zç]+\s+de\s+\d{4})\s*-\s*(Present|o momento|[a-zç]+\s+de\s+\d{4})/i;
    let currentJob = null;

    for (let i = 0; i < expLines.length; i++) {
        const line = expLines[i];
        if (dateRegex.test(line)) {
            if (currentJob) {
                currentJob.description = currentJob.description.trim();
                profile.experience.push(currentJob);
            }
            const dateObj = parseDateRange(line);
            currentJob = {
                company: expLines[i - 2] || "",
                role: expLines[i - 1] || "",
                ...dateObj,
                location: "",
                description: ""
            };
            if (expLines[i + 1] && expLines[i + 1].includes("(") && /\d/.test(expLines[i + 1])) {
                currentJob.location = expLines[i + 2] || "";
                i += 2;
            } else {
                currentJob.location = expLines[i + 1] || "";
                i += 1;
            }
        } else {
            if (currentJob && !line.toLowerCase().includes("page")) currentJob.description += line + "\n";
        }
    }
    if (currentJob) {
        currentJob.description = currentJob.description.trim();
        profile.experience.push(currentJob);
    }

    // Formação
    const eduLines = extractSection("Formação acadêmica", ["Page", "Certifications", "Publications"]);
    if (eduLines) {
        for (let i = 0; i < eduLines.length; i++) {
            let line = eduLines[i].trim();
            if (line.startsWith("·") || line.startsWith("(") || /^\W*\(?[a-zç]+\s+de\s+\d{4}/i.test(line)) {
                let fullDate = line;
                let offset = 0;
                while (!fullDate.includes(")") && (i + offset + 1) < eduLines.length) {
                    offset++;
                    fullDate += " " + eduLines[i + offset].trim();
                }
                const cleanDateStr = fullDate.replace(/[·()]/g, "").trim();
                const dateObj = parseDateRange(cleanDateStr);
                const degree = eduLines[i - 1] || "";
                const school = eduLines[i - 2] || "";

                if (degree && school) {
                    profile.education.push({ school: school.trim(), degree: degree.trim(), ...dateObj });
                }
                i += offset;
            }
        }
    }

    return profile;
}