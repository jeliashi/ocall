import 'package:flutter/material.dart';
import 'package:ocall_app/auth_gate.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // appBar: AppBar(
      //   automaticallyImplyLeading: false,
      // ),
      body: Container(
        alignment: Alignment.center,
        decoration: const BoxDecoration(
            image: DecorationImage(
                invertColors: true,
                image: AssetImage("login.jpg"),
                fit: BoxFit.cover)),
        child: Column(
          children: [
            Text('Welcome!',
                style: Theme.of(context).textTheme.displayLarge,
                selectionColor: Theme.of(context).primaryColor),
          ],
        ),
      ),
    );
  }
}
