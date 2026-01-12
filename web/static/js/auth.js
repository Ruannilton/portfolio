/**
 * auth.js - Funções de autenticação (login/cadastro)
 */

function showTab(tab) {
    const loginTab = document.getElementById('tab-login');
    const signupTab = document.getElementById('tab-signup');
    const loginForm = document.getElementById('form-login');
    const signupForm = document.getElementById('form-signup');
    const response = document.getElementById('response');
    
    response.innerHTML = '';
    
    if (tab === 'login') {
        loginTab.classList.add('text-blue-600', 'border-blue-500');
        loginTab.classList.remove('text-gray-500', 'border-transparent');
        loginForm.classList.remove('hidden');
       
        signupTab.classList.remove('text-blue-600', 'border-blue-500');
        signupTab.classList.add('text-gray-500', 'border-transparent');
        signupForm.classList.add('hidden');
    } else {
        signupTab.classList.add('text-blue-600', 'border-blue-500');
        signupTab.classList.remove('text-gray-500', 'border-transparent');
        signupForm.classList.remove('hidden');
        
        loginTab.classList.remove('text-blue-600', 'border-blue-500');
        loginTab.classList.add('text-gray-500', 'border-transparent');
        loginForm.classList.add('hidden');
    }
}

function handleAuthResponse(event) {
    const xhr = event.detail.xhr;
    const response = document.getElementById('response');
    
    if (xhr.status >= 200 && xhr.status < 300) {
        try {
            const data = JSON.parse(xhr.responseText);
            if (data.access_token) {
                // Salva token em cookie (será setado pelo servidor também)
                document.cookie = `access_token=${data.access_token}; path=/; max-age=${data.expires_in}; SameSite=Lax`;
                
                // Mostra sucesso e redireciona
                response.innerHTML = '<div class="bg-green-100 text-green-700 p-3 rounded-lg">Login realizado com sucesso! Redirecionando...</div>';
                
                setTimeout(() => {
                    window.location.href = '/app';
                }, 500);
            } else {
                // Registro sem auto-login (fallback)
                response.innerHTML = '<div class="bg-green-100 text-green-700 p-3 rounded-lg">Conta criada! Redirecionando...</div>';
                setTimeout(() => {
                    window.location.href = '/app';
                }, 500);
            }
        } catch (e) {
            // Resposta não é JSON (pode ser registro bem-sucedido sem body)
            if (xhr.status === 201) {
                response.innerHTML = '<div class="bg-green-100 text-green-700 p-3 rounded-lg">Conta criada! Redirecionando...</div>';
                setTimeout(() => {
                    window.location.href = '/app';
                }, 500);
            }
        }
    } else {
        // Erro
        response.innerHTML = `<div class="bg-red-100 text-red-700 p-3 rounded-lg">${xhr.responseText || 'Erro ao processar requisição'}</div>`;
    }
}
