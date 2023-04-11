import 'dart:convert';

import 'package:firebase_auth/firebase_auth.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class ProfileForm extends StatefulWidget {
  @override
  _ProfileFormState createState() => _ProfileFormState();
}

class _ProfileFormState extends State<ProfileForm> {
  String _profileType = '';
  String _profileName = '';
  bool _showLocationField = false;
  bool _isLoading = false;

  final _formKey = GlobalKey<FormState>();

  void _handleRadioValueChanged(String? value) {
    setState(() {
      if (value == null) {
        _profileType == "";
      } else {
        _profileType = value;
      }
      if (value == 'venue') {
        _showLocationField = true;
      } else {
        _showLocationField = false;
      }
    });
  }

  void _handleProfileNameChanged(String? value) {
    setState(() {
      if (value == null) {
        _profileName = "";
      } else {
        _profileName = value;
      }
    });
  }

  Future<String?> _handleSubmit() async {
    setState(() {
      _isLoading = true;
    });
    if (_formKey.currentState!.validate()) {
      String? authToken =
          await FirebaseAuth.instance.currentUser?.getIdToken(false);
      if (authToken == null) {
        return "Experienced error in acquiring authentication.\nPlease sign out and back in.";
        // ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
        //     content: Text(
        //         "Experienced error in acquiring authentication.\nPlease sign out and back in.")));
        // setState(() {
        //   _isLoading = false;
        // });
        // return;
      }

      Map<String, dynamic> requestBody = {
        'name': _profileName,
        'type': _profileType,
      };
      // if (_profileType == 'venue') {
      //   requestBody['location'] = location;
      // }
      try {
        http.Response response = await http.post(
          Uri.parse('https://your-api-url.com/profiles'),
          headers: {
            'Authorization': 'Bearer $authToken',
            'Content-Type': 'application/json',
          },
          body: jsonEncode(requestBody),
        );

        // Check the response status code and handle the response accordingly
        if (response.statusCode == 201) {
          // Success
          print('Profile created successfully!');
          // Navigator.of(context).pushReplacement(MaterialPageRoute(
          //     builder: (BuildContext context) => const HomeScreen()));
          return "";
        } else {
          // Error
          // print('Error creating profile: ${response.body}');
          // ScaffoldMessenger.of(context).showSnackBar(SnackBar(
          //     content: Text(
          //         "There was a server error.\nPlease try again later\n\nDetails: ${response.statusCode}: ${response.body}")));
          return "There was a server error.\nPlease try again later\n\nDetails: ${response.statusCode}: ${response.body}";
        }
      } catch (error) {
        // ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        //     content: Text(
        //         "There was an internal error.\nPlease try again later\n\nDetails: $error")));
        return "There was an internal error.\nPlease try again later\n\nDetails: $error";
      } finally {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Create Profile'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('What type of profile is this?'),
              Row(
                children: [
                  Radio<String>(
                    value: 'venue',
                    groupValue: _profileType,
                    onChanged: _handleRadioValueChanged,
                  ),
                  const Text('Venue'),
                  Radio<String>(
                    value: 'performer',
                    groupValue: _profileType,
                    onChanged: _handleRadioValueChanged,
                  ),
                  const Text('Performer'),
                  Radio<String>(
                    value: 'producer',
                    groupValue: _profileType,
                    onChanged: _handleRadioValueChanged,
                  ),
                  const Text('Producer'),
                ],
              ),
              if (_profileType.isNotEmpty)
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('What is the name of your $_profileType profile?'),
                    TextFormField(
                      onChanged: _handleProfileNameChanged,
                      validator: (value) {
                        if (value!.isEmpty) {
                          return 'Please enter a name for your profile';
                        }
                        return null;
                      },
                    ),
                  ],
                ),
              if (_showLocationField)
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('Where is your venue located?'),
                    TextFormField(
                      validator: (value) {
                        if (value!.isEmpty) {
                          return 'Please enter the location of your venue';
                        }
                        return null;
                      },
                    ),
                  ],
                ),
              const SizedBox(height: 16.0),
              ElevatedButton(
                onPressed: _isLoading
                    ? null
                    : () async {
                        var result = await _handleSubmit();
                        if ((result == "") || (result == null)) {
                          Navigator.pop(context);
                        } else {
                          ScaffoldMessenger.of(context)
                              .showSnackBar(SnackBar(content: Text(result)));
                        }
                      },
                child: _isLoading
                    ? const CircularProgressIndicator()
                    : const Text("Submit"),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
