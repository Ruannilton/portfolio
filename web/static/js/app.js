/**
 * app.js - Funções globais da aplicação
 */

// Navegação
function redirectToSearch() {
    window.location.href = '/search';
}

function login() {
    window.location.href = '/login';
}

function logout() {
    fetch('/auth/logout', { method: 'GET' })
        .finally(() => {
            window.location.href = '/login';
        });
}

// Controle de seções
function showSection(sectionId) {
    const sections = ['portfolio', 'conversations', 'stats'];
    sections.forEach(id => {
        const section = document.getElementById(id);
        if (section) {
            if (id === sectionId) {
                section.classList.remove('hidden');
            } else {
                section.classList.add('hidden');
            }
        }
    });
}
