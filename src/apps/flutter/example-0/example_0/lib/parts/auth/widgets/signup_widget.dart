import 'package:example_0/ui_kit/ui_kit_part.dart';
import 'package:flutter/material.dart';

class SignupWidget extends StatefulWidget {
  const SignupWidget({super.key});

  @override
  State<SignupWidget> createState() => _SignupWidgetState();
}

class _SignupWidgetState extends State<SignupWidget> {

  final _loginController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        HeaderTextField(
          headerText: "Login",
          controller: _loginController,          
          ),
        const SizedBox(
          height: 16,
        ),
        HeaderTextField(
          headerText: "Password",
          controller: _passwordController,          
        ),
        const SizedBox(
          height: 16,
        ),
        ElevatedButton(onPressed: 
          (){
            print("loginField: ${_loginController.text} passwordField: ${_passwordController.text} ");
          }, 
          child: const Text("Signup")
        )
      ],
    );
  }
}