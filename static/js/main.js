$(()=>{
    // formulario de registro
    $('.modal').modal();
    $("#form_registro").submit((e)=>{
        e.preventDefault();
        $.ajax({
            url: '/registro',
            type: 'POST',
            dataType: 'json',
            data: $('#form_registro').serialize(),
            success: (json) => {
                if (json.status === "success"){
                    $('#modal_reg .modal-content h4').html("Registrado con Éxito")
                    $('#modal_reg .modal-content p').html("Ahora se le llevara la página para iniciar sesión.")
                    $('#modal_reg .modal-footer a').attr("href","/login")
                } else {
                    $('#modal_reg .modal-content h4').html("Error")
                    $('#modal_reg .modal-content p').html(json.content)
                }
                $('#modal_reg').modal('open');
            }
        });
    });

    // formulario de login
    $("#form_login").submit((e) => {
        e.preventDefault();
        $.ajax({
            url: '/login',
            type: 'POST',
            dataType: 'json',
            data: $('#form_login').serialize(),
            success: (json) => {
                if (json.status === "success") {
                    location.href = "/";
                } else {
                    $('#modal_reg .modal-content h4').html("Error")
                    $('#modal_reg .modal-content p').html(json.content)
                }
                $('#modal_reg').modal('open');
            }
        });
    });
})