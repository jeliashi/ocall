import 'dart:convert';
import 'dart:io';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:flutterfire_ui/auth.dart';
import 'package:http/http.dart' as http;
import 'package:ocall_app/create_profile.dart';

class Profile {
  final String name;
  final String type;

  Profile({required this.name, required this.type});

  factory Profile.fromJson(Map<String, dynamic> json) {
    return Profile(
      name: json['name'],
      type: json['type'],
    );
  }
}

class HomeScreen extends StatefulWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  HomeScreenState createState() => HomeScreenState();
}

class HomeScreenState extends State<HomeScreen> {
  late Future<List<Profile>> _futureProfiles;
  @override
  void initState() {
    super.initState();
    _futureProfiles = _fetchProfiles();
  }

  Future<List<Profile>> _fetchProfiles() async {
    String? authToken =
        await FirebaseAuth.instance.currentUser?.getIdToken(false);
    if (authToken == null) {
      throw AuthFailed(Exception());
    }
    try {
      http.Response response = await http.get(
        Uri.parse("http://localhost:8000/profiles"),
        headers: {
          'Authorization': 'Bearer $authToken',
          'Content-Type': 'application/json',
        },
      );
      if (response.statusCode == 200) {
        final List<dynamic> data = jsonDecode(response.body);
        return data.map((json) => Profile.fromJson(json)).toList();
      } else {
        throw HttpException(
            "Server Error. Please try again later.\nDetails ${response.statusCode}: ${response.body}");
      }
    } catch (err) {
      rethrow;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
      ),
      body: Container(
        alignment: Alignment.center,
        // decoration: const BoxDecoration(
        //     image: DecorationImage(
        //         invertColors: true,
        //         image: AssetImage("login.jpg"),
        //         fit: BoxFit.cover)),
        child: Column(
          children: [
            Text('Welcome!',
                style: Theme.of(context).textTheme.displayLarge,
                selectionColor: Theme.of(context).primaryColor),
            FutureBuilder<List<Profile>>(
              future: _futureProfiles,
              builder: (context, snapshot) {
                if (snapshot.hasData) {
                  final List<Profile> profiles = snapshot.data!;
                  return ListView.builder(
                    itemCount: profiles.length,
                    itemBuilder: (context, index) {
                      final profile = profiles[index];
                      return ListTile(
                        title: Text(profile.name),
                        subtitle: Text(profile.type),
                      );
                    },
                  );
                } else if (snapshot.hasError) {
                  return Center(child: Text('${snapshot.error}'));
                } else {
                  return const Center(child: CircularProgressIndicator());
                }
              },
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(builder: (context) => ProfileForm()),
                );
              },
              child: const Text('Create Profile'),
            )
          ],
        ),
      ),
    );
  }
}
