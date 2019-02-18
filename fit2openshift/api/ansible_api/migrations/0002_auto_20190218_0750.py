# Generated by Django 2.1.2 on 2019-02-18 07:50

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('openshift_api', '0006_auto_20190218_0750'),
        ('ansible_api', '0001_initial'),
    ]

    operations = [
        migrations.AlterField(
            model_name='Group',
            name='hosts',
            field=models.ManyToManyField(related_name='groups', to='ansible_api.Host'),
        ),
    ]
