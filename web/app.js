$(document).ready(function () {
	const API_BASE_URL = "http://localhost:8080/api/v1";
	let autoRefreshInterval = null;

	// Отправка выражения
	$("#submitBtn").click(function () {
		const expression = $("#expressionInput").val().trim();

		if (!expression) {
			alert("Введите выражение");
			return;
		}

		$.ajax({
			url: `${API_BASE_URL}/calculate`,
			method: "POST",
			contentType: "application/json",
			data: JSON.stringify({ expression }),
			success: function () {
				$("#expressionInput").val("");
				loadExpressions();
			},
			error: function (xhr) {
				alert(`Ошибка: ${xhr.statusText}`);
			},
		});
	});

	// Загрузка списка выражений
	function loadExpressions() {
		$.ajax({
			url: `${API_BASE_URL}/expressions`,
			method: "GET",
			success: function (data) {
				renderExpressions(data.expressions);
			},
			error: function (err) {
				alert("Не удалось загрузить выражения" + err);
			},
		});
	}

	// Отрисовка выражений
	function renderExpressions(expressions) {
		console.log(expressions)
		const items = expressions
			.map(
				(expr) => `
          <div class="expression-item">
              <div>
                  <strong>${expr.expression}</strong><br>
                  <small>ID: ${expr.id}</small>
              </div>
              <div class="status-${expr.status.toLowerCase()}">
                  ${expr.status}${expr.result ? `: ${parseFloat(expr.result)}` : ""}
              </div>
          </div>
      `
			)
			.join("");

		$("#expressionsList").html(items);
	}

	// Обработка автообновления
	$("#autoRefresh").change(function () {
		if ($(this).is(":checked")) {
			autoRefreshInterval = setInterval(loadExpressions, 2000);
			loadExpressions();
		} else {
			clearInterval(autoRefreshInterval);
		}
	});

	// Инициализация
	$("#refreshBtn").click(loadExpressions);
	//loadExpressions();
});
