$(document).ready(function() {
	const API_BASE_URL = 'http://localhost:8080/api/v1';
	let autoRefreshInterval = null;

	$('#expressionInput').keypress(function (e) {
		if (e.which === 13) {
			$('#submitBtn').click();
		}
	});

	$('#submitBtn').click(function() {
		loadExpressions();
	})

	$('#submitBtn').click(function() {
		const expression = $('#expressionInput').val().trim();

		if (!expression) {
			showError('Введите выражение');
			return;
		}

		$(this).prop('disabled', true).html('<i class="fas fa-spinner fa-spin"></i> Отправка...');

		$.ajax({
			url: `${API_BASE_URL}/calculate`,
			method: 'POST',
			contentType: 'application/json',
			data: JSON.stringify({ expression }),
			success: function() {
				$('#expressionInput').val('');
				loadExpressions();
			},
			error: function(xhr) {
				showError(xhr.responseJSON?.error || 'Ошибка сервера');
			},
			complete: function() {
				$('#submitBtn').prop('disabled', false).html('<i class="fas fa-paper-plane"></i> Отправить');
			}
		});
	});

	function loadExpressions() {
		$('#expressionsList').html('<div class="loading"><i class="fas fa-spinner fa-spin"></i> Загрузка...</div>');

		$.ajax({
			url: `${API_BASE_URL}/expressions`,
			method: 'GET',
			success: function(data) {
				renderExpressions(data.expressions);
			},
			error: function() {
				showError('Не удалось загрузить данные');
			}
		});
	}

	function renderExpressions(expressions) {
		const items = expressions.map(expr => `
            <div class="expression-item status-${expr.status.toLowerCase()}">
                <div class="expression-info">
                    <div>${expr.expression}</div>
                    <div class="expression-id">ID: ${expr.id}</div>
                </div>
                <div class="expression-status">
                    ${statusIcon(expr.status)} ${expr.status}${expr.result ? `: ${expr.result}` : ''}
                </div>
            </div>
        `).join('');

		$('#expressionsList').html(items || '<div class="empty-state">Нет активных вычислений</div>');
	}

	function statusIcon(status) {
		const icons = {
			pending: 'fa-clock',
			processing: 'fa-sync-alt fa-spin',
			completed: 'fa-check-circle',
			error: 'fa-exclamation-circle'
		};
		return `<i class="fas ${icons[status.toLowerCase()]}"></i>`;
	}

	function showError(message) {
		const errorEl = $('<div class="error-message"><i class="fas fa-exclamation-triangle"></i> ' + message + '</div>');
		$('.container').prepend(errorEl);
		setTimeout(() => errorEl.remove(), 3000);
	}

	$('#autoRefresh').change(function() {
		if ($(this).is(':checked')) {
			autoRefreshInterval = setInterval(loadExpressions, 1000);
			loadExpressions();
		} else {
			clearInterval(autoRefreshInterval);
		}
	});

	//loadExpressions();
});