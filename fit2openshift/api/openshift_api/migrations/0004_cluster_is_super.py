# Generated by Django 2.1.2 on 2019-02-18 03:20

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('openshift_api', '0003_deployexecution_operation'),
    ]

    operations = [
        migrations.AddField(
            model_name='cluster',
            name='is_super',
            field=models.BooleanField(default=False),
        ),
    ]
